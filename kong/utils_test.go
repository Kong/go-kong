package kong

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

const StatsDSchema = `{
		"name" : "statsd",
		"fields" : [
			{
				"config" : {
					"type" : "record",
					"fields": [
						{
							"host": {
								"default": "localhost",
								"type": "string"
							}
						},
						{
							"port": {
								"between": [
									0,
									65535
								],
								"default": 8125,
								"type": "integer"
							}
						},
						{
							"prefix": {
								"default": "kong",
								"type": "string"
							}
						},
						{
							"metrics": {
								"default": [
									{
										"name": "request_count",
										"sample_rate": 1,
										"stat_type": "counter"
									},
									{
										"name": "latency",
										"stat_type": "timer"
									},
									{
										"name": "request_size",
										"stat_type": "timer"
									},
									{
										"name": "status_count",
										"sample_rate": 1,
										"stat_type": "counter"
									},
									{
										"name": "response_size",
										"stat_type": "timer"
									},
									{
										"consumer_identifier": "custom_id",
										"name": "unique_users",
										"stat_type": "set"
									},
									{
										"consumer_identifier": "custom_id",
										"name": "request_per_user",
										"sample_rate": 1,
										"stat_type": "counter"
									},
									{
										"name": "upstream_latency",
										"stat_type": "timer"
									},
									{
										"name": "kong_latency",
										"stat_type": "timer"
									},
									{
										"consumer_identifier": "custom_id",
										"name": "status_count_per_user",
										"sample_rate": 1,
										"stat_type": "counter"
									}
								],
								"elements": {
									"entity_checks": [
										{
											"conditional": {
												"if_field": "name",
												"if_match": {
													"eq": "unique_users"
												},
												"then_field": "stat_type",
												"then_match": {
													"eq": "set"
												}
											}
										},
										{
											"conditional": {
												"if_field": "stat_type",
												"if_match": {
													"one_of": [
														"counter",
														"gauge"
													]
												},
												"then_field": "sample_rate",
												"then_match": {
													"required": true
												}
											}
										},
										{
											"conditional": {
												"if_field": "name",
												"if_match": {
													"one_of": [
														"status_count_per_user",
														"request_per_user",
														"unique_users"
													]
												},
												"then_field": "consumer_identifier",
												"then_match": {
													"required": true
												}
											}
										},
										{
											"conditional": {
												"if_field": "name",
												"if_match": {
													"one_of": [
														"status_count",
														"status_count_per_user",
														"request_per_user"
													]
												},
												"then_field": "stat_type",
												"then_match": {
													"eq": "counter"
												}
											}
										}
									],
									"fields": [
										{
											"name": {
												"one_of": [
													"kong_latency",
													"latency",
													"request_count",
													"request_per_user",
													"request_size",
													"response_size",
													"status_count",
													"status_count_per_user",
													"unique_users",
													"upstream_latency"
												],
												"required": true,
												"type": "string"
											}
										},
										{
											"stat_type": {
												"one_of": [
													"counter",
													"gauge",
													"histogram",
													"meter",
													"set",
													"timer"
												],
												"required": true,
												"type": "string"
											}
										},
										{
											"sample_rate": {
												"gt": 0,
												"type": "number"
											}
										},
										{
											"consumer_identifier": {
												"one_of": [
													"consumer_id",
													"custom_id",
													"username"
												],
												"type": "string"
											}
										}
									],
									"type": "record"
								},
								"type": "array"
							}
						}
					]
				}
			}
		]
	}`

// TestSchemaSetType (
const TestSchemaSetType = `{
	"fields": [{
		"config": {
			"type": "record",
			"fields": [{
				"bootstrap_servers": {
					"default": [
						{
							"host": "127.0.0.1",
							"port": 42
						}
					],
					"elements": {
						"fields": [{
								"host": {
									"required": true,
									"type": "string"
								}
							},
							{
								"port": {
									"between": [
										0,
										65535
									],
									"required": true,
									"type": "integer",
									"default": 42
								}
							}
						],
						"type": "record"
					},
					"type": "set"
				}
			}]
		}
	}]
}`

const AcmeSchema = `{
    "fields": [
        {
            "consumer": {
                "type": "foreign",
                "eq": null,
                "reference": "consumers"
            }
        },
        {
            "service": {
                "type": "foreign",
                "eq": null,
                "reference": "services"
            }
        },
        {
            "route": {
                "type": "foreign",
                "eq": null,
                "reference": "routes"
            }
        },
        {
            "protocols": {
                "elements": {
                    "one_of": [
                        "grpc",
                        "grpcs",
                        "http",
                        "https"
                    ],
                    "type": "string"
                },
                "required": true,
                "default": [
                    "grpc",
                    "grpcs",
                    "http",
                    "https"
                ],
                "type": "set"
            }
        },
        {
            "config": {
                "required": true,
                "fields": [
                    {
                        "account_email": {
                            "referenceable": true,
                            "type": "string",
                            "match": "%w*%p*@+%w*%.?%w*",
                            "required": true,
                            "encrypted": true
                        }
                    },
                    {
                        "account_key": {
                            "required": false,
                            "fields": [
                                {
                                    "key_id": {
                                        "type": "string",
                                        "required": true
                                    }
                                },
                                {
                                    "key_set": {
                                        "type": "string"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    },
                    {
                        "api_uri": {
                            "default": "https://acme-v02.api.letsencrypt.org/directory",
                            "type": "string"
                        }
                    },
                    {
                        "tos_accepted": {
                            "default": false,
                            "type": "boolean"
                        }
                    },
                    {
                        "eab_kid": {
                            "referenceable": true,
                            "encrypted": true,
                            "type": "string"
                        }
                    },
                    {
                        "eab_hmac_key": {
                            "referenceable": true,
                            "encrypted": true,
                            "type": "string"
                        }
                    },
                    {
                        "cert_type": {
                            "one_of": [
                                "rsa",
                                "ecc"
                            ],
                            "default": "rsa",
                            "type": "string"
                        }
                    },
                    {
                        "rsa_key_size": {
                            "one_of": [
                                2048,
                                3072,
                                4096
                            ],
                            "default": 4096,
                            "type": "number"
                        }
                    },
                    {
                        "renew_threshold_days": {
                            "default": 14,
                            "type": "number"
                        }
                    },
                    {
                        "domains": {
                            "elements": {
                                "match_any": {
                                    "patterns": [
                                        "^%*%.",
                                        "%.%*$",
                                        "^[^*]*$"
                                    ],
                                    "err": "invalid wildcard: must be placed at leftmost or rightmost label"
                                },
                                "match_all": [
                                    {
                                        "pattern": "^[^*]*%*?[^*]*$",
                                        "err": "invalid wildcard: must have at most one wildcard"
                                    }
                                ],
                                "type": "string"
                            },
                            "type": "array"
                        }
                    },
                    {
                        "allow_any_domain": {
                            "default": false,
                            "type": "boolean"
                        }
                    },
                    {
                        "fail_backoff_minutes": {
                            "default": 5,
                            "type": "number"
                        }
                    },
                    {
                        "storage": {
                            "one_of": [
                                "kong",
                                "shm",
                                "redis",
                                "consul",
                                "vault"
                            ],
                            "default": "shm",
                            "type": "string"
                        }
                    },
                    {
                        "storage_config": {
                            "required": true,
                            "fields": [
                                {
                                    "shm": {
                                        "required": true,
                                        "fields": [
                                            {
                                                "shm_name": {
                                                    "default": "kong",
                                                    "type": "string"
                                                }
                                            }
                                        ],
                                        "type": "record"
                                    }
                                },
                                {
                                    "kong": {
                                        "required": true,
                                        "fields": [],
                                        "type": "record"
                                    }
                                },
                                {
                                    "redis": {
                                        "required": true,
                                        "fields": [
                                            {
                                                "host": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "port": {
                                                    "between": [
                                                        0,
                                                        65535
                                                    ],
                                                    "type": "integer"
                                                }
                                            },
                                            {
                                                "database": {
                                                    "type": "number"
                                                }
                                            },
                                            {
                                                "auth": {
                                                    "referenceable": true,
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "ssl": {
                                                    "required": true,
                                                    "default": false,
                                                    "type": "boolean"
                                                }
                                            },
                                            {
                                                "ssl_verify": {
                                                    "required": true,
                                                    "default": false,
                                                    "type": "boolean"
                                                }
                                            },
                                            {
                                                "ssl_server_name": {
                                                    "type": "string",
                                                    "required": false
                                                }
                                            },
                                            {
                                                "namespace": {
                                                    "len_min": 0,
                                                    "required": true,
                                                    "default": "",
                                                    "type": "string"
                                                }
                                            }
                                        ],
                                        "type": "record"
                                    }
                                },
                                {
                                    "consul": {
                                        "required": true,
                                        "fields": [
                                            {
                                                "https": {
                                                    "default": false,
                                                    "type": "boolean"
                                                }
                                            },
                                            {
                                                "host": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "port": {
                                                    "between": [
                                                        0,
                                                        65535
                                                    ],
                                                    "type": "integer"
                                                }
                                            },
                                            {
                                                "kv_path": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "timeout": {
                                                    "type": "number"
                                                }
                                            },
                                            {
                                                "token": {
                                                    "referenceable": true,
                                                    "type": "string"
                                                }
                                            }
                                        ],
                                        "type": "record"
                                    }
                                },
                                {
                                    "vault": {
                                        "required": true,
                                        "fields": [
                                            {
                                                "https": {
                                                    "default": false,
                                                    "type": "boolean"
                                                }
                                            },
                                            {
                                                "host": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "port": {
                                                    "between": [
                                                        0,
                                                        65535
                                                    ],
                                                    "type": "integer"
                                                }
                                            },
                                            {
                                                "kv_path": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "timeout": {
                                                    "type": "number"
                                                }
                                            },
                                            {
                                                "token": {
                                                    "referenceable": true,
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "tls_verify": {
                                                    "default": true,
                                                    "type": "boolean"
                                                }
                                            },
                                            {
                                                "tls_server_name": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "auth_method": {
                                                    "one_of": [
                                                        "token",
                                                        "kubernetes"
                                                    ],
                                                    "default": "token",
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "auth_path": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "auth_role": {
                                                    "type": "string"
                                                }
                                            },
                                            {
                                                "jwt_path": {
                                                    "type": "string"
                                                }
                                            }
                                        ],
                                        "type": "record"
                                    }
                                }
                            ],
                            "type": "record"
                        }
                    },
                    {
                        "preferred_chain": {
                            "type": "string"
                        }
                    },
                    {
                        "enable_ipv4_common_name": {
                            "default": true,
                            "type": "boolean"
                        }
                    }
                ],
                "type": "record"
            }
        }
    ]
}`

const defaultRecordSchema = `{
	"fields": [
			{
					"config": {
							"fields": [
									{
											"endpoint": {
													"description": "endpoint",
													"referenceable": true,
													"required": true,
													"type": "string"
											}
									},
									{
											"queue": {
													"default": {
															"max_batch_size": 200
													},
													"fields": [
															{
																	"max_batch_size": {
																			"between": [
																					1,
																					1000000
																			],
																			"default": 1,
																			"type": "integer"
																	}
															},
															{
																	"max_coalescing_delay": {
																			"between": [
																					0,
																					3600
																			],
																			"default": 1,
																			"type": "number"
																	}
															}
													],
													"required": true,
													"type": "record"
											}
									},
									{
											"propagation": {
													"default": {
															"default_format": "w3c"
													},
													"fields": [
															{
																	"extract": {
																			"elements": {
																					"one_of": [
																							"w3c",
																							"b3"
																					],
																					"type": "string"
																			},
																			"type": "array"
																	}
															},
															{
																	"default_format": {
																			"one_of": [
																					"w3c",
																					"b3"
																			],
																			"required": true,
																			"type": "string"
																	}
															}
													],
													"required": true,
													"type": "record"
											}
									}
							],
							"required": true,
							"type": "record"
					}
			}
	]
}`

func TestStringArrayToString(t *testing.T) {
	assert := assert.New(t)

	arr := StringSlice("foo", "bar")
	s := stringArrayToString(arr)
	assert.Equal("[ foo, bar ]", s)

	arr = StringSlice("foo")
	s = stringArrayToString(arr)
	assert.Equal("[ foo ]", s)

	assert.Equal("nil", stringArrayToString(nil))
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	s := String("foo")
	assert.Equal("foo", *s)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	b := Bool(true)
	assert.True(*b)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)

	i := Int(42)
	assert.Equal(42, *i)
}

func TestStringSlice(t *testing.T) {
	assert := assert.New(t)

	arrp := StringSlice()
	assert.Empty(arrp)

	arrp = StringSlice("foo", "bar")
	assert.Len(arrp, 2)
	assert.Equal("foo", *arrp[0])
	assert.Equal("bar", *arrp[1])
}

func TestFixVersion(t *testing.T) {
	tests := []struct {
		version         string
		expectedVersion string
		isEnterprise    bool
	}{
		{
			version:         "0.14.1",
			expectedVersion: "0.14.1",
		},
		{
			version:         "0.14.2rc",
			expectedVersion: "0.14.2",
		},
		{
			version:         "0.14.2rc1",
			expectedVersion: "0.14.2",
		},
		{
			version:         "0.14.2preview",
			expectedVersion: "0.14.2",
		},
		{
			version:         "0.14.2preview1",
			expectedVersion: "0.14.2",
		},
		{
			version:         "0.33-enterprise-edition",
			expectedVersion: "0.33.0",
			isEnterprise:    true,
		},
		{
			version:         "0.33-1-enterprise-edition",
			expectedVersion: "0.33.1",
			isEnterprise:    true,
		},
		{
			version:         "1.3.0.0-enterprise-edition-lite",
			expectedVersion: "1.3.0.0",
			isEnterprise:    true,
		},
		{
			version:         "3.0.0.0",
			expectedVersion: "3.0.0.0",
			isEnterprise:    true,
		},
		{
			version:         "3.0.0.0-enterprise-edition",
			expectedVersion: "3.0.0.0",
			isEnterprise:    true,
		},
	}
	for _, test := range tests {
		v, err := ParseSemanticVersion(test.version)
		if err != nil {
			t.Errorf("error converting %s: %v", test.version, err)
		} else if v.String() != test.expectedVersion {
			t.Errorf("converting %s, expecting %s, getting %s", test.version, test.expectedVersion, v.String())
		}
		assert.Equal(t, test.isEnterprise, v.IsKongGatewayEnterprise())
	}

	invalidVersions := []string{
		"",
		"0-1-1",
	}
	for _, inputVersion := range invalidVersions {
		_, err := ParseSemanticVersion(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}

func Test_requestWithHeaders(t *testing.T) {
	type args struct {
		req     *http.Request
		headers http.Header
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "returns request as is if no headers are set",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo": []string{"bar", "baz"},
				},
			},
		},
		{
			name: "returns request with headers added",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo":      []string{"bar", "baz"},
					"Password": []string{"my-secret-key"},
					"Key-With": []string{"multiple", "values"},
				},
			},
		},
		{
			name: "returns nil when input request is nil",
			args: args{
				req: nil,
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestWithHeaders(tt.args.req, tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestWithHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillRoutesDefaults(T *testing.T) {
	SkipWhenKongRouterFlavor(T, Expressions)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		route    *Route
		expected *Route
	}{
		{
			name: "fills defaults for all fields except paths, leaves name unchanged",
			route: &Route{
				Name:  String("r1"),
				Paths: []*string{String("/r1")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("http"), String("https")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
		{
			name: "fills defaults for all fields except paths and protocols, leaves name unchanged",
			route: &Route{
				Name:      String("r1"),
				Paths:     []*string{String("/r1")},
				Protocols: []*string{String("grpc")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("grpc")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
		{
			name: "boolean default values don't overwrite existing fields if set",
			route: &Route{
				Name:         String("r1"),
				Paths:        []*string{String("/r1")},
				Protocols:    []*string{String("grpc")},
				StripPath:    Bool(false),
				PreserveHost: Bool(true),
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(true),
				Protocols:               []*string{String("grpc")},
				RegexPriority:           Int(0),
				StripPath:               Bool(false),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			r := tc.route
			fullSchema, err := client.Schemas.Get(defaultCtx, "routes")
			require.NoError(T, err)
			assert.NotNil(fullSchema)
			require.NoError(t, FillEntityDefaults(r, fullSchema))
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(
				Route{},
				"RequestBuffering", "ResponseBuffering", "PathHandling",
			)
			if diff := cmp.Diff(r, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillServiceDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		service  *Service
		expected *Service
	}{
		{
			name: "fills defaults for all fields, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(80),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name, tags and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
				Tags: []*string{String("tag1"), String("tag2")},
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
				Tags:           []*string{String("tag1"), String("tag2")},
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			s := tc.service
			fullSchema, err := client.Schemas.Get(defaultCtx, "services")
			require.NoError(T, err)
			assert.NotNil(fullSchema)
			require.NoError(t, FillEntityDefaults(s, fullSchema))
			opt := []cmp.Option{
				cmpopts.IgnoreFields(Service{}, "Enabled"),
			}
			if diff := cmp.Diff(s, tc.expected, opt...); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillTargetDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		target   *Target
		expected *Target
	}{
		{
			name:   "fills default for weight",
			target: &Target{},
			expected: &Target{
				Weight: Int(100),
			},
		},
		{
			name: "set zero-value",
			target: &Target{
				Weight: Int(0),
			},
			expected: &Target{
				Weight: Int(0),
			},
		},
	}

	kongVersion := GetVersionForTesting(T)
	hasFailoverVersionRange := MustNewRange(">=3.12.0")
	shouldContainFailover := hasFailoverVersionRange(kongVersion)

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			target := tc.target
			fullSchema, err := client.Schemas.Get(defaultCtx, "targets")
			require.NoError(t, err)
			require.NotNil(t, fullSchema)
			require.NoError(t, FillEntityDefaults(target, fullSchema))

			// Gateway 3.12 added a new Failover field to targets
			// which has a default of False
			if shouldContainFailover {
				tc.expected.Failover = Bool(false)
			}

			if diff := cmp.Diff(target, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillUpstreamsDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		upstream *Upstream
		expected *Upstream
	}{
		{
			name: "fills defaults for all fields, leaves name unchanged",
			upstream: &Upstream{
				Name: String("upstream1"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("round-robin"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("none"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
		{
			name: "fills defaults for all fields except algorithm and hash_on, leaves name unchanged",
			upstream: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				HashOn:    String("ip"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("ip"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
	}
	// Kong Enterprise added `StickySessionsCookiePath` field and assigned default values since 3.11.
	// We need to add the default value of `StickySessionsCookiePath` when Kong version >= 3.11.0.
	kongVersion := GetVersionForTesting(T)
	hasStickySessionsCookiePathVersionRange := MustNewRange(">=3.11.0")

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			u := tc.upstream
			fullSchema, err := client.Schemas.Get(defaultCtx, "upstreams")
			require.NoError(T, err)
			assert.NotNil(fullSchema)
			require.NoError(t, FillEntityDefaults(u, fullSchema))
			// Ignore fields to make tests pass despite small differences across releases.
			opts := []cmp.Option{
				cmpopts.IgnoreFields(Healthcheck{}, "Threshold"),
				cmpopts.IgnoreFields(Upstream{}, "UseSrvName"),
			}

			// Add the default value of `throttling` to the expected configuration of plugin.
			if hasStickySessionsCookiePathVersionRange(kongVersion) {
				tc.expected.StickySessionsCookiePath = String("/")
			}

			if diff := cmp.Diff(u, tc.expected, opts...); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestUpstreamStickySessionsFields(t *testing.T) {
	tests := []struct {
		name     string
		upstream *Upstream
		expected *Upstream
	}{
		{
			name: "sticky sessions cookie field is preserved",
			upstream: &Upstream{
				Name:                 String("test-upstream"),
				StickySessionsCookie: String("session_id"),
			},
			expected: &Upstream{
				Name:                 String("test-upstream"),
				StickySessionsCookie: String("session_id"),
			},
		},
		{
			name: "sticky sessions cookie path field is preserved",
			upstream: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookiePath: String("/api"),
			},
			expected: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookiePath: String("/api"),
			},
		},
		{
			name: "both sticky sessions fields are preserved",
			upstream: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookie:     String("session_id"),
				StickySessionsCookiePath: String("/api"),
			},
			expected: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookie:     String("session_id"),
				StickySessionsCookiePath: String("/api"),
			},
		},
		{
			name: "sticky sessions fields work with nil values",
			upstream: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookie:     nil,
				StickySessionsCookiePath: nil,
			},
			expected: &Upstream{
				Name:                     String("test-upstream"),
				StickySessionsCookie:     nil,
				StickySessionsCookiePath: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test that the fields are correctly set and preserved
			assert.Equal(t, tc.expected.Name, tc.upstream.Name)
			assert.Equal(t, tc.expected.StickySessionsCookie, tc.upstream.StickySessionsCookie)
			assert.Equal(t, tc.expected.StickySessionsCookiePath, tc.upstream.StickySessionsCookiePath)
		})
	}
}

func TestUpstreamStickySessionsJSONSerialization(t *testing.T) {
	upstream := &Upstream{
		Name:                     String("test-upstream"),
		StickySessionsCookie:     String("session_id"),
		StickySessionsCookiePath: String("/api"),
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(upstream)
	require.NoError(t, err)

	// Verify the JSON contains the expected fields
	assert.Contains(t, string(jsonData), `"sticky_sessions_cookie":"session_id"`)
	assert.Contains(t, string(jsonData), `"sticky_sessions_cookie_path":"/api"`)

	// Test JSON unmarshaling
	var unmarshaledUpstream Upstream
	err = json.Unmarshal(jsonData, &unmarshaledUpstream)
	require.NoError(t, err)

	// Verify the fields were correctly unmarshaled
	assert.Equal(t, "test-upstream", *unmarshaledUpstream.Name)
	assert.Equal(t, "session_id", *unmarshaledUpstream.StickySessionsCookie)
	assert.Equal(t, "/api", *unmarshaledUpstream.StickySessionsCookiePath)
}

func getJSONSchemaFromFile(t *testing.T, filename string) Schema {
	jsonFile, err := os.Open(filename)
	require.NoError(t, err)
	defer jsonFile.Close()
	var schema Schema
	require.NoError(t, json.NewDecoder(jsonFile).Decode(&schema))
	return schema
}

func TestFillUpstreamsDefaultsFromJSONSchema(t *testing.T) {
	// load upstream JSON schema from local file.
	schema := getJSONSchemaFromFile(t, "testdata/upstreamJSONSchema.json")

	tests := []struct {
		name     string
		upstream *Upstream
		expected *Upstream
	}{
		{
			name: "fills defaults for all fields, leaves name unchanged",
			upstream: &Upstream{
				Name: String("upstream1"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("round-robin"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("none"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
		{
			name: "fills defaults for all fields except algorithm and hash_on, leaves name unchanged",
			upstream: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				HashOn:    String("ip"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("ip"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.upstream
			require.NoError(t, FillEntityDefaults(u, schema))
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(Healthcheck{}, "Threshold")
			if diff := cmp.Diff(u, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillServicesDefaultsFromJSONSchema(t *testing.T) {
	// load service JSON schema from local file.
	schema := getJSONSchemaFromFile(t, "testdata/serviceJSONSchema.json")

	tests := []struct {
		name     string
		service  *Service
		expected *Service
	}{
		{
			name: "fills defaults for all fields, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(80),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name, tags and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
				Tags: []*string{String("tag1"), String("tag2")},
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
				Tags:           []*string{String("tag1"), String("tag2")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.service
			require.NoError(t, FillEntityDefaults(s, schema))
			opt := []cmp.Option{
				cmpopts.IgnoreFields(Service{}, "Enabled"),
			}
			if diff := cmp.Diff(s, tc.expected, opt...); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillRoutesDefaultsFromJSONSchema(t *testing.T) {
	// load route JSON schema from local file.
	schema := getJSONSchemaFromFile(t, "testdata/routeJSONSchema.json")

	tests := []struct {
		name     string
		route    *Route
		expected *Route
	}{
		{
			name: "fills defaults for all fields except paths, leaves name unchanged",
			route: &Route{
				Name:  String("r1"),
				Paths: []*string{String("/r1")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("http"), String("https")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
		{
			name: "fills defaults for all fields except paths and protocols, leaves name unchanged",
			route: &Route{
				Name:      String("r1"),
				Paths:     []*string{String("/r1")},
				Protocols: []*string{String("grpc")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("grpc")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
		{
			name: "boolean default values don't overwrite existing fields if set",
			route: &Route{
				Name:         String("r1"),
				Paths:        []*string{String("/r1")},
				Protocols:    []*string{String("grpc")},
				StripPath:    Bool(false),
				PreserveHost: Bool(true),
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(true),
				Protocols:               []*string{String("grpc")},
				RegexPriority:           Int(0),
				StripPath:               Bool(false),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.route
			require.NoError(t, FillEntityDefaults(r, schema))
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(
				Route{},
				"RequestBuffering", "ResponseBuffering", "PathHandling",
			)
			if diff := cmp.Diff(r, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestFillTargetDefaultsFromJSONSchema(t *testing.T) {
	// load route JSON schema from local file.
	schema := getJSONSchemaFromFile(t, "testdata/targetJSONSchema.json")

	tests := []struct {
		name     string
		target   *Target
		expected *Target
	}{
		{
			name:   "fills default for weight",
			target: &Target{},
			expected: &Target{
				Weight: Int(100),
			},
		},
		{
			name: "set zero-value",
			target: &Target{
				Weight: Int(0),
			},
			expected: &Target{
				Weight: Int(0),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			target := tc.target
			require.NoError(t, FillEntityDefaults(target, schema))
			require.Equal(t, tc.expected, target)
		})
	}
}

func TestHTTPClientWithHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	assert.NotPanics(t,
		func() {
			client := HTTPClientWithHeaders(&http.Client{}, nil)
			assert.NotNil(t, client)
		},
		"creating Kong's HTTP client using default/uninitialized http.Client shouldn't panic",
	)

	assert.NotPanics(t,
		func() {
			client := HTTPClientWithHeaders(nil, nil)
			assert.NotNil(t, client)
		},
		"creating Kong's HTTP client using nil http.Client shouldn't panic",
	)
}

func TestFillConsumerGroupPluginDefaults(T *testing.T) {
	RunWhenEnterprise(T, ">=2.7.0", RequiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		plugin   *ConsumerGroupPlugin
		expected *ConsumerGroupPlugin
	}{
		{
			name:   "fills default for consumer_group_plugins",
			plugin: &ConsumerGroupPlugin{},
			expected: &ConsumerGroupPlugin{
				Config: Configuration{
					"window_type":            "sliding",
					"retry_after_jitter_max": float64(0),
				},
			},
		},
		{
			name: "fills default only for unset retry_after_jitter_max field",
			plugin: &ConsumerGroupPlugin{
				Config: Configuration{
					"window_type": "fixed",
				},
			},
			expected: &ConsumerGroupPlugin{
				Config: Configuration{
					"window_type":            "fixed",
					"retry_after_jitter_max": float64(0),
				},
			},
		},
		{
			name: "fills default only for unset window_type field",
			plugin: &ConsumerGroupPlugin{
				Config: Configuration{
					"retry_after_jitter_max": float64(10),
				},
			},
			expected: &ConsumerGroupPlugin{
				Config: Configuration{
					"window_type":            "sliding",
					"retry_after_jitter_max": float64(10),
				},
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin
			fullSchema, err := client.Schemas.Get(defaultCtx, "consumer_group_plugins")
			require.NoError(T, err)
			assert.NotNil(fullSchema)
			require.NoError(t, FillEntityDefaults(plugin, fullSchema))
			if diff := cmp.Diff(plugin, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

const fillConfigRecordTestSchema = `{
	"fields": {
		"config": {
			"type": "record",
			"fields": [
				{
					"enabled": {
						"type": "boolean",
						"default": true,
						"required": true
					}
				},
				{
					"mappings": {
						"required": false,
						"type": "array",
						"elements": {
							"type": "record",
							"fields": [
								{
									"name": {
										"type": "string",
										"required": false
									}
								},
								{
									"nationality": {
										"type": "string",
										"required": false
									}
								}
							]
						}
					}
				},
				{
					"empty_record": {
						"type": "record",
						"required": true,
						"fields": []
					}
				}
			]
		}
	}
}
`

func Test_fillConfigRecord(t *testing.T) {
	tests := []struct {
		name     string
		schema   gjson.Result
		config   Configuration
		expected Configuration
	}{
		{
			name:   "fills defaults for all missing fields",
			schema: gjson.Parse(fillConfigRecordTestSchema),
			config: Configuration{
				"mappings": []any{
					map[string]any{
						"nationality": "Ethiopian",
					},
				},
			},
			expected: Configuration{
				"enabled": true,
				"mappings": []any{
					Configuration{
						"name":        nil,
						"nationality": "Ethiopian",
					},
				},
			},
		},
		{
			name:   "handle empty array as nil for a record field",
			schema: gjson.Parse(fillConfigRecordTestSchema),
			config: Configuration{
				"mappings": []any{
					map[string]any{
						"nationality": "Ethiopian",
					},
				},
				"empty_record": map[string]any{},
			},
			expected: Configuration{
				"enabled": true,
				"mappings": []any{
					Configuration{
						"name":        nil,
						"nationality": "Ethiopian",
					},
				},
				"empty_record": map[string]any{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configSchema, err := getConfigSchema(tc.schema)
			require.NoError(t, err)
			config := fillConfigRecord(configSchema, tc.config, nil, FillRecordOptions{
				FillDefaults: true,
				FillAuto:     true,
			})
			require.NotNil(t, config)
			if diff := cmp.Diff(config, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

const fillConfigRecordTestSchemaWithShorthandFields = `{
	"fields": {
		"config": {
			"type": "record",
			"shorthand_fields": [
				{
					"redis_port": {
						"translate_backwards": [
							"redis",
							"port"
						],
						"type": "integer"
					}
				},
				{
					"redis_host": {
						"translate_backwards": [
							"redis",
							"host"
						],
						"type": "string"
					}
				}
			],
			"fields": [
				{
					"enabled": {
						"type": "boolean",
						"default": true,
						"required": true
					}
				},
				{
					"mappings": {
						"required": false,
						"type": "array",
						"elements": {
							"type": "record",
							"fields": [
								{
									"name": {
										"type": "string",
										"required": false
									}
								},
								{
									"nationality": {
										"type": "string",
										"required": false
									}
								}
							]
						}
					}
				},
				{
					"empty_record": {
						"type": "record",
						"required": true,
						"fields": []
					}
				},
				{
					"redis": {
						"required": true,
						"description": "Redis configuration",
						"type": "record",
						"fields": [
							{
								"host": {
									"type": "string"
								}
							},
							{
								"port": {
									"default": 6379,
									"type": "integer"
								}
							}
						]
					}
				}
			]
		}
	}
}
`

const fillConfigRecordTestSchemaWithAutoFields = `{
	"fields": {
		"config": {
			"type": "record",
			"fields": [
				{
					"default_string": {
						"type": "string",
						"default": "abc"
					}
				},
				{
					"auto_string_1": {
						"type": "string",
						"auto": true
					}
				},
				{
					"auto_string_2": {
						"type": "string",
						"auto": true
					}
				},
				{
					"auto_string_3": {
						"type": "string",
						"auto": true
					}
				}
			]
		}
	}
}
`

const fillConfigRecordTestSchemaWithRecord = `{
	"fields": {
		"config": {
			"type": "record",
			"fields": [
				{
					"some_record": {
							"required": true,
              "fields": [
                {
                  "some_field": {
                    "default": "kong",
                    "type": "string"
                  }
                }
              ],
              "type": "record"
          }
				},
				{
					"some_other_record": {
              "fields": [
                {
                  "some_field": {
                    "type": "string"
                  }
                }
              ],
              "type": "record"
          }
				},
				{
					"string_1": {
						"type": "string",
					}
				}
			]
		}
	}
}
`

func Test_fillConfigRecord_shorthand_fields(t *testing.T) {
	tests := []struct {
		name     string
		schema   gjson.Result
		config   Configuration
		expected Configuration
	}{
		{
			name:   "fills defaults for all missing fields",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithShorthandFields),
			config: Configuration{
				"mappings": []any{
					map[string]any{
						"nationality": "Ethiopian",
					},
				},
			},
			expected: Configuration{
				"enabled": true,
				"mappings": []any{
					Configuration{
						"name":        nil,
						"nationality": "Ethiopian",
					},
				},
				"redis": map[string]interface{}{
					"host": nil,
					"port": float64(6379),
				},
			},
		},
		{
			name:   "backfills nested fields if shorthand field values are changed",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithShorthandFields),
			config: Configuration{
				"redis_host": "localhost",
				"redis_port": float64(8000),
			},
			expected: Configuration{
				"enabled":    true,
				"mappings":   nil,
				"redis_port": float64(8000),
				"redis_host": "localhost",
			},
		},
		{
			name:   "backfills nested fields if shorthand field values are changed and respects nil value (over default)",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithShorthandFields),
			config: Configuration{
				"redis_host": "localhost-custom-1",
				"redis_port": nil,
			},
			expected: Configuration{
				"enabled":    true,
				"mappings":   nil,
				"redis_port": nil, // new field redis.port has defined default value but the given redis_port is respected
				"redis_host": "localhost-custom-1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configSchema, err := getConfigSchema(tc.schema)
			require.NoError(t, err)
			config := fillConfigRecord(configSchema, tc.config, nil, FillRecordOptions{
				FillDefaults: true,
				FillAuto:     true,
			})
			require.NotNil(t, config)
			if diff := cmp.Diff(config, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_fillConfigRecord_defaults_only(t *testing.T) {
	tests := []struct {
		name     string
		schema   gjson.Result
		config   Configuration
		expected Configuration
	}{
		{
			name:   "fills defaults with opts to fill defaults only",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithAutoFields),
			config: Configuration{
				"auto_string_3": "789",
			},
			expected: Configuration{
				"default_string": "abc",
				"auto_string_3":  "789",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configSchema, err := getConfigSchema(tc.schema)
			require.NoError(t, err)
			config := fillConfigRecord(configSchema, tc.config, nil, FillRecordOptions{
				FillDefaults: true,
				FillAuto:     false,
			})
			require.NotNil(t, config)
			if diff := cmp.Diff(config, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_fillConfigRecord_auto_only(t *testing.T) {
	tests := []struct {
		name     string
		schema   gjson.Result
		config   Configuration
		expected Configuration
	}{
		{
			name:   "fills defaults with opts to fill auto only",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithAutoFields),
			config: Configuration{
				"auto_string_3": "789",
			},
			expected: Configuration{
				"auto_string_1": nil,
				"auto_string_2": nil,
				"auto_string_3": "789",
				// defalt_string missing
			},
		},
		{
			name:   "not passing record field leaves field unset",
			schema: gjson.Parse(fillConfigRecordTestSchemaWithRecord),
			config: Configuration{
				// some_record missing
				"some_other_record": map[string]any{}, // explicitly set to empty record
				"string_1":          "abc",
			},
			expected: Configuration{
				// some_record was not filled
				"some_other_record": map[string]any{}, // empty record remained unchanged
				"string_1":          "abc",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configSchema, err := getConfigSchema(tc.schema)
			require.NoError(t, err)
			config := fillConfigRecord(configSchema, tc.config, nil, FillRecordOptions{
				FillDefaults: false,
				FillAuto:     true,
			})
			require.NotNil(t, config)
			if diff := cmp.Diff(config, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaults(t *testing.T) {
	defaultMetrics := []any{
		map[string]any{
			"name":        "request_count",
			"sample_rate": float64(1),
			"stat_type":   "counter",
		},
		map[string]any{
			"name":      "latency",
			"stat_type": "timer",
		},
		map[string]any{
			"name":      "request_size",
			"stat_type": "timer",
		},
		map[string]any{
			"name":        "status_count",
			"sample_rate": float64(1),
			"stat_type":   "counter",
		},
		map[string]any{
			"name":      "response_size",
			"stat_type": "timer",
		},
		map[string]any{
			"consumer_identifier": "custom_id",
			"name":                "unique_users",
			"stat_type":           "set",
		},
		map[string]any{
			"consumer_identifier": "custom_id",
			"name":                "request_per_user",
			"sample_rate":         float64(1),
			"stat_type":           "counter",
		},
		map[string]any{
			"name":      "upstream_latency",
			"stat_type": "timer",
		},
		map[string]any{
			"name":      "kong_latency",
			"stat_type": "timer",
		},
		map[string]any{
			"consumer_identifier": "custom_id",
			"name":                "status_count_per_user",
			"sample_rate":         float64(1),
			"stat_type":           "counter",
		},
	}
	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "fills defaults for all missing fields",
			plugin: &Plugin{
				Config: Configuration{
					"metrics": []interface{}{
						map[string]interface{}{
							"name":      "response_size",
							"stat_type": "histogram",
						},
					},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"host":   "localhost",
					"port":   float64(8125),
					"prefix": "kong",
					"metrics": []interface{}{
						Configuration{
							"name":                "response_size",
							"stat_type":           "histogram",
							"consumer_identifier": nil,
							"sample_rate":         nil,
						},
					},
				},
			},
		},
		{
			name: "fills defaults when empty array of records in config",
			plugin: &Plugin{
				Config: Configuration{
					"metrics": []any{},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"host":    "localhost",
					"port":    float64(8125),
					"prefix":  "kong",
					"metrics": defaultMetrics,
				},
			},
		},
		{
			name: "fills defaults when nil array of records in config",
			plugin: &Plugin{
				Config: Configuration{
					"metrics": nil,
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"host":    "localhost",
					"port":    float64(8125),
					"prefix":  "kong",
					"metrics": nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin

			var fullSchema map[string]interface{}
			err := json.Unmarshal([]byte(StatsDSchema), &fullSchema)
			require.NoError(t, err)
			require.NotNil(t, fullSchema)

			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin,
				"Protocols", "Enabled",
			)
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaults_RequestTransformer(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "fills defaults for all missing fields",
			plugin: &Plugin{
				Config: Configuration{
					"add": map[string]interface{}{
						"body": []any{},
						"headers": []any{
							"Knative-Serving-Namespace:e3ffeafd-b5fe-4f34-b2e4-af6f3d9fb417",
							"Knative-Serving-Revision:helloworld-go-00001",
						},
						"querystring": []any{},
					},
					"append": map[string]interface{}{
						"body":        []any{},
						"headers":     []any{},
						"querystring": []any{},
					},
					"http_method": nil,
					"enabled":     true,
					"id":          "0beef60e-e7e3-40f8-ac47-f6a10b931cee",
					"name":        "request-transformer",
					"protocols": []any{
						"grpc",
						"grpcs",
						"http",
						"https",
					},
					"service": map[string]interface{}{
						"id":   "63295454-c41e-447e-bce5-d6b488f3866e",
						"name": "30bc1240-ad84-4760-a469-763bce524191.helloworld-go-00001.80",
					},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"add": map[string]any{
						"body": []any{},
						"headers": []any{
							"Knative-Serving-Namespace:e3ffeafd-b5fe-4f34-b2e4-af6f3d9fb417",
							"Knative-Serving-Revision:helloworld-go-00001",
						},
						"querystring": []any{},
					},
					"append": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": []interface{}{},
					},
					"remove":      map[string]any{"body": []any{}, "headers": []any{}, "querystring": []any{}},
					"rename":      map[string]any{"body": []any{}, "headers": []any{}, "querystring": []any{}},
					"replace":     map[string]any{"body": []any{}, "headers": []any{}, "querystring": []any{}, "uri": nil},
					"http_method": nil,
					"enabled":     true,
					"id":          "0beef60e-e7e3-40f8-ac47-f6a10b931cee",
					"name":        "request-transformer",
					"protocols": []any{
						"grpc",
						"grpcs",
						"http",
						"https",
					},
					"service": map[string]interface{}{
						"id":   "63295454-c41e-447e-bce5-d6b488f3866e",
						"name": "30bc1240-ad84-4760-a469-763bce524191.helloworld-go-00001.80",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin
			fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/request-transformer")
			require.NoError(t, err)
			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaults_SetType(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "does not fill defaults when provided",
			plugin: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []interface{}{
						map[string]interface{}{
							"host": "192.168.2.100",
							"port": float64(3500),
						},
						map[string]any{
							"host": "192.168.2.101",
							"port": float64(3500),
						},
					},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []interface{}{
						Configuration{
							"host": "192.168.2.100",
							"port": float64(3500),
						},
						Configuration{
							"host": "192.168.2.101",
							"port": float64(3500),
						},
					},
				},
			},
		},
		{
			name: "fills defaults for all missing fields",
			plugin: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []interface{}{
						map[string]interface{}{
							"host": "127.0.0.1",
						},
					},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []interface{}{
						Configuration{
							"host": "127.0.0.1",
							"port": float64(42),
						},
					},
				},
			},
		},
		{
			name: "fills defaults when empty set of records in config",
			plugin: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []any{},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"bootstrap_servers": []any{
						map[string]any{
							"host": "127.0.0.1",
							"port": float64(42),
						},
					},
				},
			},
		},
		{
			name: "fills defaults when nil set of records in config",
			plugin: &Plugin{
				Config: Configuration{
					"bootstrap_servers": nil,
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"bootstrap_servers": nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin

			var fullSchema map[string]interface{}
			err := json.Unmarshal([]byte(TestSchemaSetType), &fullSchema)
			require.NoError(t, err)
			require.NotNil(t, fullSchema)

			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin,
				"Protocols", "Enabled",
			)
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaults_Acme(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "fills defaults for all missing fields",
			plugin: &Plugin{
				Config: Configuration{},
			},
			expected: &Plugin{
				Config: Configuration{
					"account_email":           nil,
					"account_key":             nil,
					"allow_any_domain":        bool(false),
					"api_uri":                 string("https://acme-v02.api.letsencrypt.org/directory"),
					"cert_type":               string("rsa"),
					"domains":                 nil,
					"eab_hmac_key":            nil,
					"eab_kid":                 nil,
					"enable_ipv4_common_name": bool(true),
					"fail_backoff_minutes":    float64(5),
					"preferred_chain":         nil,
					"renew_threshold_days":    float64(14),
					"rsa_key_size":            float64(4096),
					"storage":                 string("shm"),
					"storage_config": map[string]any{
						"consul": map[string]any{
							"host":    nil,
							"https":   bool(false),
							"kv_path": nil,
							"port":    nil,
							"timeout": nil,
							"token":   nil,
						},
						"redis": map[string]any{
							"auth":            nil,
							"database":        nil,
							"host":            nil,
							"namespace":       string(""),
							"port":            nil,
							"ssl":             bool(false),
							"ssl_server_name": nil,
							"ssl_verify":      bool(false),
						},
						"shm": map[string]any{"shm_name": string("kong")},
						"vault": map[string]any{
							"auth_method":     string("token"),
							"auth_path":       nil,
							"auth_role":       nil,
							"host":            nil,
							"https":           bool(false),
							"jwt_path":        nil,
							"kv_path":         nil,
							"port":            nil,
							"timeout":         nil,
							"tls_server_name": nil,
							"tls_verify":      bool(true),
							"token":           nil,
						},
					},
					"tos_accepted": bool(false),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin
			var fullSchema map[string]interface{}
			err := json.Unmarshal([]byte(AcmeSchema), &fullSchema)

			require.NoError(t, err)
			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaults_DefaultRecord(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "record defaults take precedence over fields defaults",
			plugin: &Plugin{
				Config: Configuration{
					"endpoint": "http://test.test:4317",
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"endpoint": "http://test.test:4317",
					"propagation": map[string]interface{}{
						"default_format": string("w3c"), // from record defaults
						"extract":        nil,           // from field defaults
					},
					"queue": map[string]interface{}{
						"max_batch_size":       float64(200), // from record defaults
						"max_coalescing_delay": float64(1),   // from field defaults
					},
				},
			},
		},
		{
			name: "configured values take precedence over record defaults",
			plugin: &Plugin{
				Config: Configuration{
					"endpoint": "http://test.test:4317",
					"propagation": map[string]interface{}{
						"default_format": "b3",
					},
					"queue": map[string]interface{}{
						"max_batch_size": 123,
					},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"endpoint": "http://test.test:4317",
					"propagation": map[string]interface{}{
						"default_format": string("b3"),
						"extract":        nil, // from field defaults
					},
					"queue": map[string]interface{}{
						"max_batch_size":       float64(123),
						"max_coalescing_delay": float64(1), // from field defaults
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin
			var fullSchema map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(defaultRecordSchema), &fullSchema))
			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

const NonEmptyDefaultArrayFieldSchema = `{
    "fields": [
        {
            "protocols": {
                "default": [
                    "grpc",
                    "grpcs",
                    "http",
                    "https"
                ],
                "elements": {
                    "len_min": 1,
                    "one_of": [
                        "grpc",
                        "grpcs",
                        "http",
                        "https"
                    ],
                    "required": true,
                    "type": "string"
                },
                "required": true,
                "type": "set"
            }
        },
        {
            "config": {
                "fields": [
                    {
                        "issuer": {
                            "required": true,
                            "type": "string"
                        }
                    },
                    {
                        "login_tokens": {
                            "default": [
                                "id_token"
                            ],
                            "elements": {
                                "one_of": [
                                    "id_token",
                                    "access_token",
                                    "refresh_token",
                                    "tokens",
                                    "introspection"
                                ],
                                "type": "string"
                            },
                            "required": false,
                            "type": "array"
                        }
                    }
                ],
                "type": "record"
            }
        }
    ]
}`

func Test_FillPluginsDefaults_NonEmptyDefaultArrayField(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "not setting login_tokens should be overwritten by default value",
			plugin: &Plugin{
				Config: Configuration{
					"issuer": "https://accounts.google.com",
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"issuer":       "https://accounts.google.com",
					"login_tokens": []any{"id_token"},
				},
			},
		},
		{
			name: "setting empty array for login_tokens should not be overwritten by default value",
			plugin: &Plugin{
				Config: Configuration{
					"issuer":       "https://accounts.google.com",
					"login_tokens": []any{},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"issuer":       "https://accounts.google.com",
					"login_tokens": []any{},
				},
			},
		},
		{
			name: "setting non-empty login_tokens should not be overwritten by default value",
			plugin: &Plugin{
				Config: Configuration{
					"issuer":       "https://accounts.google.com",
					"login_tokens": []any{"access_token", "refresh_token"},
				},
			},
			expected: &Plugin{
				Config: Configuration{
					"issuer":       "https://accounts.google.com",
					"login_tokens": []any{"access_token", "refresh_token"},
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			plugin := tc.plugin
			var fullSchema map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(NonEmptyDefaultArrayFieldSchema), &fullSchema))

			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_ClearUnmatchingDeprecationsSimple(t *testing.T) {
	RunWhenKong(t, ">=3.8.0")
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting")
	require.NoError(t, err)
	require.NotNil(t, fullSchema)

	tests := []struct {
		name              string
		newPlugin         *Plugin
		oldPlugin         *Plugin
		expectedOldPlugin Configuration
	}{
		{
			name: "when new object contains only old (deprecated) fields",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis_host": "localhost",
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"host": "localhost",
					},
					"redis_host": "localhost",
				},
			},
			expectedOldPlugin: Configuration{
				"redis_host": "localhost",
			},
		},
		{
			name: "when new object contains only new fields (non-deprecated)",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"host": "localhost",
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"host": "localhost",
					},
					"redis_host": "localhost",
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"host": "localhost",
				},
			},
		},
		{
			name: "when new object contains both new and old fields",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"host": "localhost",
					},
					"redis_host": "localhost",
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"host": "localhost",
					},
					"redis_host": "localhost",
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"host": "localhost",
				},
				"redis_host": "localhost",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ClearUnmatchingDeprecations(tc.newPlugin, tc.oldPlugin, fullSchema))
			if diff := cmp.Diff(tc.oldPlugin.Config, tc.expectedOldPlugin); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_ClearUnmatchingDeprecationsAdvanced(t *testing.T) {
	RunWhenEnterprise(t, ">=3.8.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, fullSchema)

	tests := []struct {
		name              string
		newPlugin         *Plugin
		oldPlugin         *Plugin
		expectedOldPlugin Configuration
	}{
		{
			name: "when new object contains only old (deprecated) fields",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
				},
			},
		},
		{
			name: "when new object contains only new fields",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
				},
			},
		},
		{
			name: "when new object contains old field but the new ones are split into multiple separate fields",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"timeout": 2000,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"timeout": 2000,
				},
			},
		},
		{
			name: "when new object contains new field that is split into multiple fields but there was only one old field",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
		{
			name: "when both complete new and old configuration is sent",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
					"timeout":         2000,
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
		{
			name: "when both complete new and old configuration is sent but their values are nil",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": nil,
						"cluster_nodes":     nil,
						"timeout":           nil,
						"connect_timeout":   nil,
						"send_timeout":      nil,
						"read_timeout":      nil,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": nil,
						"cluster_nodes":     nil,
						"timeout":           nil,
						"connect_timeout":   nil,
						"send_timeout":      nil,
						"read_timeout":      nil,
					},
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": nil,
					"cluster_nodes":     nil,
					"timeout":           nil,
					"connect_timeout":   nil,
					"send_timeout":      nil,
					"read_timeout":      nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ClearUnmatchingDeprecations(tc.newPlugin, tc.oldPlugin, fullSchema))
			if diff := cmp.Diff(tc.oldPlugin.Config, tc.expectedOldPlugin); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_ClearUnmatchingDeprecationsWhenSchemaIsWrong(t *testing.T) {
	tests := []struct {
		name   string
		schema map[string]interface{}
	}{
		// These test cases are rather theoretical since the schema is a JSON extracted from Kong /schemas endpoint
		{
			name: "when schema is not json serializble",
			schema: map[string]interface{}{
				"some other field": math.Inf(1),
			},
		},
		{
			name: "when schema is wrong - i.e. does not have {fields: [ {config: {fields: []}} ]} structure",
			schema: map[string]interface{}{
				"some other field": 4,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Error(t, ClearUnmatchingDeprecations(nil, nil, tc.schema))
		})
	}
}

func Test_ClearUnmatchingDeprecationsWhenNotUpdateEvent(t *testing.T) {
	RunWhenEnterprise(t, ">=3.8.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, fullSchema)

	tests := []struct {
		name                     string
		newPlugin                *Plugin
		oldPlugin                *Plugin
		expectedNewPluginCleared Configuration
		expectedOldPlugin        Configuration
	}{
		{
			name: "when only new configuration is sent (CREATE event)",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			oldPlugin:         nil,
			expectedOldPlugin: nil,
			expectedNewPluginCleared: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
					"timeout":         2000,
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
		{
			name:      "when only old configuration is sent (DELETE event)",
			newPlugin: nil,
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedNewPluginCleared: nil,
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
					"timeout":         2000,
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ClearUnmatchingDeprecations(tc.newPlugin, tc.oldPlugin, fullSchema))
			if tc.expectedNewPluginCleared != nil {
				if diff := cmp.Diff(tc.newPlugin.Config, tc.expectedNewPluginCleared); diff != "" {
					t.Errorf("unexpected diff:\n%s", diff)
				}
			}

			if tc.expectedOldPlugin != nil {
				if diff := cmp.Diff(tc.oldPlugin.Config, tc.expectedOldPlugin); diff != "" {
					t.Errorf("unexpected diff:\n%s", diff)
				}
			}
		})
	}
}

func Test_ClearUnmatchingDeprecationsWhenNewConfigIsSetAsNil(t *testing.T) {
	RunWhenEnterprise(t, ">=3.8.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, fullSchema)

	tests := []struct {
		name                     string
		newPlugin                *Plugin
		oldPlugin                *Plugin
		expectedNewPluginCleared Configuration
		expectedOldPlugin        Configuration
	}{
		{
			name: "when only old configuration is sent but the new one was filled with nil",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes":     nil,
						"timeout":           2000,
						"connect_timeout":   nil,
						"send_timeout":      nil,
						"read_timeout":      nil,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedNewPluginCleared: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"timeout":           2000,
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"timeout":           2000,
				},
			},
		},
		{
			name: "when both new and old configuration is sent and their values differ - (should not change configurations)",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 9379},
							{"ip": "127.0.0.1", "port": 9380},
							{"ip": "127.0.0.1", "port": 9381},
						},
						"timeout":         2000,
						"connect_timeout": 3001,
						"send_timeout":    3002,
						"read_timeout":    3003,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedNewPluginCleared: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 9379},
						{"ip": "127.0.0.1", "port": 9380},
						{"ip": "127.0.0.1", "port": 9381},
					},
					"timeout":         2000,
					"connect_timeout": 3001,
					"send_timeout":    3002,
					"read_timeout":    3003,
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
					"timeout":         2000,
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ClearUnmatchingDeprecations(tc.newPlugin, tc.oldPlugin, fullSchema))
			if diff := cmp.Diff(tc.newPlugin.Config, tc.expectedNewPluginCleared); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}

			if diff := cmp.Diff(tc.oldPlugin.Config, tc.expectedOldPlugin); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_ClearUnmatchingDeprecationsWhenNewConfigHasDefaults(t *testing.T) {
	RunWhenEnterprise(t, ">=3.8.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	fullSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, fullSchema)

	tests := []struct {
		name                     string
		newPlugin                *Plugin
		oldPlugin                *Plugin
		expectedNewPluginCleared Configuration
		expectedOldPlugin        Configuration
	}{
		{
			name: "when only old configuration is sent but the new one was filled with nil",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes":     nil,
						"timeout":           2000,
						"connect_timeout":   nil,
						"send_timeout":      nil,
						"read_timeout":      nil,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedNewPluginCleared: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"timeout":           2000,
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"timeout":           2000,
				},
			},
		},
		{
			name: "when both new and old configuration is sent and their values differ - (should not change configurations)",
			newPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 9379},
							{"ip": "127.0.0.1", "port": 9380},
							{"ip": "127.0.0.1", "port": 9381},
						},
						"timeout":         2000,
						"connect_timeout": 3001,
						"send_timeout":    3002,
						"read_timeout":    3003,
					},
				},
			},
			oldPlugin: &Plugin{
				Config: Configuration{
					"redis": map[string]interface{}{
						"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
						"cluster_nodes": []map[string]interface{}{
							{"ip": "127.0.0.1", "port": 6379},
							{"ip": "127.0.0.1", "port": 6380},
							{"ip": "127.0.0.1", "port": 6381},
						},
						"timeout":         2000,
						"connect_timeout": 2000,
						"send_timeout":    2000,
						"read_timeout":    2000,
					},
				},
			},
			expectedNewPluginCleared: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 9379},
						{"ip": "127.0.0.1", "port": 9380},
						{"ip": "127.0.0.1", "port": 9381},
					},
					"timeout":         2000,
					"connect_timeout": 3001,
					"send_timeout":    3002,
					"read_timeout":    3003,
				},
			},
			expectedOldPlugin: Configuration{
				"redis": map[string]interface{}{
					"cluster_addresses": []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
					"cluster_nodes": []map[string]interface{}{
						{"ip": "127.0.0.1", "port": 6379},
						{"ip": "127.0.0.1", "port": 6380},
						{"ip": "127.0.0.1", "port": 6381},
					},
					"timeout":         2000,
					"connect_timeout": 2000,
					"send_timeout":    2000,
					"read_timeout":    2000,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, ClearUnmatchingDeprecations(tc.newPlugin, tc.oldPlugin, fullSchema))
			if diff := cmp.Diff(tc.newPlugin.Config, tc.expectedNewPluginCleared); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}

			if diff := cmp.Diff(tc.oldPlugin.Config, tc.expectedOldPlugin); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPartialDefaults(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	partialSchema, err := client.Schemas.Get(defaultCtx, "partials/redis-ee")
	require.NoError(t, err)
	require.NotNil(t, partialSchema)

	tests := []struct {
		name            string
		partial         *Partial
		expectedPartial *Partial
		schema          Schema
		wantErr         bool
	}{
		{
			name: "empty partial config gets filled with defaults",
			partial: &Partial{
				Config: nil,
			},
			expectedPartial: &Partial{
				Config: Configuration{
					"cluster_max_redirections": float64(5),
					"cluster_nodes":            nil,
					"connect_timeout":          float64(2000),
					"connection_is_proxied":    bool(false),
					"database":                 float64(0),
					"host":                     string("127.0.0.1"),
					"keepalive_backlog":        nil,
					"keepalive_pool_size":      float64(256),
					"password":                 nil,
					"port":                     float64(6379),
					"read_timeout":             float64(2000),
					"send_timeout":             float64(2000),
					"sentinel_master":          nil,
					"sentinel_nodes":           nil,
					"sentinel_password":        nil,
					"sentinel_role":            nil,
					"sentinel_username":        nil,
					"server_name":              nil,
					"ssl":                      bool(false),
					"ssl_verify":               bool(false),
					"username":                 nil,
				},
			},
			schema: partialSchema,
		},
		{
			name: "existing partial config gets merged with defaults",
			partial: &Partial{
				Config: Configuration{
					"port":     float64(7000),
					"username": string("test-user"),
					"password": string("test-password"),
				},
			},
			expectedPartial: &Partial{
				Config: Configuration{
					"cluster_max_redirections": float64(5),
					"cluster_nodes":            nil,
					"connect_timeout":          float64(2000),
					"connection_is_proxied":    bool(false),
					"database":                 float64(0),
					"host":                     string("127.0.0.1"),
					"keepalive_backlog":        nil,
					"keepalive_pool_size":      float64(256),
					"password":                 string("test-password"),
					"port":                     float64(7000),
					"read_timeout":             float64(2000),
					"send_timeout":             float64(2000),
					"sentinel_master":          nil,
					"sentinel_nodes":           nil,
					"sentinel_password":        nil,
					"sentinel_role":            nil,
					"sentinel_username":        nil,
					"server_name":              nil,
					"ssl":                      bool(false),
					"ssl_verify":               bool(false),
					"username":                 string("test-user"),
				},
			},
			schema: partialSchema,
		},
		{
			name: "invalid schema should return error",
			partial: &Partial{
				Config: nil,
			},
			schema: Schema{
				"type": "invalid_schema",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := FillPartialDefaults(tc.partial, tc.schema)
			if tc.wantErr {
				if err == nil {
					t.Errorf("FillPartialDefaults() expected error but got none")
				}
				return
			}

			require.NoError(t, err)
			if diff := cmp.Diff(tc.partial.Config, tc.expectedPartial.Config); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginWithPartials(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	rlaPluginSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, rlaPluginSchema)

	reqPluginSchema, err := client.Schemas.Get(defaultCtx, "plugins/request-transformer")
	require.NoError(t, err)
	require.NotNil(t, reqPluginSchema)

	tests := []struct {
		name                  string
		plugin                *Plugin
		pluginSchema          map[string]interface{}
		partials              []*Partial
		expectedConfiguration Configuration
		wantErr               bool
		errString             string
	}{
		{
			name: "empty partials",
			plugin: &Plugin{
				Config: nil,
			},
			pluginSchema:          rlaPluginSchema,
			partials:              nil,
			expectedConfiguration: Configuration{},
		},
		{
			name: "plugin with single partial, path is not defined",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			expectedConfiguration: Configuration{
				"redis": Configuration{"host": string("127.0.0.1"), "port": float64(7000)},
			},
		},
		{
			name: "plugin with single partial, path is defined",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			expectedConfiguration: Configuration{
				"redis": Configuration{"host": string("127.0.0.1"), "port": float64(7000)},
			},
		},
		{
			name: "plugin that does not support partials",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: reqPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			wantErr:   true,
			errString: "schema does not contain supported_partials",
		},
		{
			name: "partial added is not supported",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: map[string]interface{}{
				"supported_partials": map[string]interface{}{},
			},
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			wantErr:   true,
			errString: "schema does not contain default partial path for partial type redis-ee",
		},
		{
			name: ">1 path found for partial added, no path defined",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: map[string]interface{}{
				"supported_partials": map[string]interface{}{
					"redis-ee": []string{"config.redis", "config.redis2"},
				},
			},
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			wantErr:   true,
			errString: ">1 supported paths found for partial type redis-ee; provide a path in config",
		},
		{
			name: ">1 path found for partial added, path is defined",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
			pluginSchema: map[string]interface{}{
				"supported_partials": map[string]interface{}{
					"redis-ee": []string{"config.redis", "config.redis2"},
				},
			},
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			expectedConfiguration: Configuration{
				"redis": Configuration{"host": string("127.0.0.1"), "port": float64(7000)},
			},
		},
		{
			name: "partial added does not exist",
			plugin: &Plugin{
				Config: nil,
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("xyz"),
					Type: String("redis-ee"),
					Config: Configuration{
						"host": string("127.0.0.1"),
						"port": float64(7000),
					},
				},
			},
			wantErr:   true,
			errString: "partial with ID abc not found",
		},
		// TODO: add a usecase for plugin that allows multiple partials
		// once that is turned on in the gateway
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := FillPluginWithPartials(tc.plugin, tc.pluginSchema, tc.partials)
			if tc.wantErr {
				if err == nil {
					t.Errorf("FillPluginWithPartials() expected error but got none")
				}
				assert.ErrorContains(t, err, tc.errString)
				return
			}

			require.NoError(t, err)
			if diff := cmp.Diff(tc.plugin.Config, tc.expectedConfiguration); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaultsWithPartials_312x(t *testing.T) {
	RunWhenEnterprise(t, ">=3.12.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	rlaPluginSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, rlaPluginSchema)

	tests := []struct {
		name           string
		plugin         *Plugin
		pluginSchema   map[string]interface{}
		partials       []*Partial
		expectedPlugin *Plugin
		wantErr        bool
		errString      string
	}{
		{
			name: "empty config, no partials present",
			plugin: &Plugin{
				Config: Configuration{},
			},
			pluginSchema: rlaPluginSchema,
			partials:     nil,
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": map[string]any{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(6379),
						"read_timeout":             float64(2000),
						"redis_proxy_type":         nil,
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"throttling":             nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
			},
		},
		{
			name: "fill plugin defaults, single partial present",
			plugin: &Plugin{
				Config: Configuration{},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
				},
			},
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"throttling":             nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
		},
		{
			name: "fill plugin defaults, single partial and path present",
			plugin: &Plugin{
				Config: Configuration{},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
				},
			},
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"throttling":             nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
		},
		{
			name: "invalid schema passed",
			plugin: &Plugin{
				ID:     String("test-plugin"),
				Config: Configuration{},
			},
			pluginSchema: map[string]interface{}{
				"type": "invalid",
			},
			partials:  nil,
			wantErr:   true,
			errString: "no 'config' field found in schema",
		},
	}
	// Kong Enterprise added `throttling` field and assigned default values since 3.12.
	// https://github.com/Kong/kong-ee/pull/12579
	// We need to add the default value of `throttling` when Kong version >= 3.12.0.
	kongVersion := GetVersionForTesting(t)
	// for nightly tests, we assume the Kong version is >= 3.13.0
	rlaHasThrottlingVersionRange := MustNewRange(">=3.13.0")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(kongVersion)

			err := FillPluginsDefaultsWithPartials(tc.plugin, tc.pluginSchema, tc.partials)
			if tc.wantErr {
				if err == nil {
					t.Errorf("FillPluginsDefaultsWithPartials expected error but got none")
				}
				assert.ErrorContains(t, err, tc.errString)
				return
			}

			require.NoError(t, err)

			// Add the default value of `throttling` to the expected configuration of plugin.
			if rlaHasThrottlingVersionRange(kongVersion) && tc.expectedPlugin != nil {
				t.Log("Add default throttling")
				tc.expectedPlugin.Config["throttling"] = nil
			}

			opts := cmpopts.IgnoreFields(*tc.plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(tc.plugin, tc.expectedPlugin, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func Test_FillPluginsDefaultsWithPartials(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0 <3.12.0", RequiredFeatures{})
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	rlaPluginSchema, err := client.Schemas.Get(defaultCtx, "plugins/rate-limiting-advanced")
	require.NoError(t, err)
	require.NotNil(t, rlaPluginSchema)

	tests := []struct {
		name           string
		plugin         *Plugin
		pluginSchema   map[string]interface{}
		partials       []*Partial
		expectedPlugin *Plugin
		wantErr        bool
		errString      string
	}{
		{
			name: "empty config, no partials present",
			plugin: &Plugin{
				Config: Configuration{},
			},
			pluginSchema: rlaPluginSchema,
			partials:     nil,
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": map[string]any{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(6379),
						"read_timeout":             float64(2000),
						"redis_proxy_type":         nil,
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
			},
		},
		{
			name: "fill plugin defaults, single partial present",
			plugin: &Plugin{
				Config: Configuration{},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
				},
			},
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
		},
		{
			name: "fill plugin defaults, single partial and path present",
			plugin: &Plugin{
				Config: Configuration{},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
			pluginSchema: rlaPluginSchema,
			partials: []*Partial{
				{
					ID:   String("abc"),
					Type: String("redis-ee"),
					Config: Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
				},
			},
			expectedPlugin: &Plugin{
				Config: Configuration{
					"compound_identifier":     nil,
					"consumer_groups":         nil,
					"dictionary_name":         string("kong_rate_limiting_counters"),
					"disable_penalty":         bool(false),
					"enforce_consumer_groups": bool(false),
					"error_code":              float64(429),
					"error_message":           string("API rate limit exceeded"),
					"header_name":             nil,
					"hide_client_headers":     bool(false),
					"identifier":              string("consumer"),
					"limit":                   nil,
					"lock_dictionary_name":    string("kong_locks"),
					"namespace":               nil,
					"path":                    nil,
					"redis": Configuration{
						"cluster_max_redirections": float64(5),
						"cluster_nodes":            nil,
						"connect_timeout":          float64(2000),
						"connection_is_proxied":    bool(false),
						"database":                 float64(0),
						"host":                     string("127.0.0.1"),
						"keepalive_backlog":        nil,
						"keepalive_pool_size":      float64(256),
						"password":                 nil,
						"port":                     float64(7000),
						"read_timeout":             float64(2000),
						"send_timeout":             float64(2000),
						"sentinel_master":          nil,
						"sentinel_nodes":           nil,
						"sentinel_password":        nil,
						"sentinel_role":            nil,
						"sentinel_username":        nil,
						"server_name":              nil,
						"ssl":                      bool(false),
						"ssl_verify":               bool(false),
						"username":                 nil,
					},
					"retry_after_jitter_max": float64(0),
					"strategy":               string("local"),
					"sync_rate":              nil,
					"window_size":            nil,
					"window_type":            string("sliding"),
				},
				Partials: []*PartialLink{
					{
						Partial: &Partial{
							ID: String("abc"),
						},
						Path: String("config.redis"),
					},
				},
			},
		},
		{
			name: "invalid schema passed",
			plugin: &Plugin{
				ID:     String("test-plugin"),
				Config: Configuration{},
			},
			pluginSchema: map[string]interface{}{
				"type": "invalid",
			},
			partials:  nil,
			wantErr:   true,
			errString: "no 'config' field found in schema",
		},
	}
	// Kong Enterprise added `throttling` field and assigned default values since 3.12.
	// https://github.com/Kong/kong-ee/pull/12579
	// We need to add the default value of `throttling` when Kong version >= 3.12.0.
	kongVersion := GetVersionForTesting(t)
	// for nightly tests, we assume the Kong version is >= 3.13.0
	rlaHasThrottlingVersionRange := MustNewRange(">=3.13.0")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(kongVersion)

			err := FillPluginsDefaultsWithPartials(tc.plugin, tc.pluginSchema, tc.partials)
			if tc.wantErr {
				if err == nil {
					t.Errorf("FillPluginsDefaultsWithPartials expected error but got none")
				}
				assert.ErrorContains(t, err, tc.errString)
				return
			}

			require.NoError(t, err)

			// Add the default value of `throttling` to the expected configuration of plugin.
			if rlaHasThrottlingVersionRange(kongVersion) && tc.expectedPlugin != nil {
				t.Log("Add default throttling")
				tc.expectedPlugin.Config["throttling"] = nil
			}

			opts := cmpopts.IgnoreFields(*tc.plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(tc.plugin, tc.expectedPlugin, opts); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}
