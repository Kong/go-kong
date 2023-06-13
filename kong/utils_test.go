package kong

import (
	"encoding/json"
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

func TestStringArrayToString(t *testing.T) {
	assert := assert.New(t)

	arr := StringSlice("foo", "bar")
	s := stringArrayToString(arr)
	assert.Equal("[ foo, bar ]", s)

	arr = StringSlice("foo")
	s = stringArrayToString(arr)
	assert.Equal("[ foo ]", s)

	assert.Equal(stringArrayToString(nil), "nil")
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	s := String("foo")
	assert.Equal("foo", *s)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	b := Bool(true)
	assert.Equal(true, *b)
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
	assert.Equal(2, len(arrp))
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
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err = FillEntityDefaults(r, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(
				Route{},
				"RequestBuffering", "ResponseBuffering", "PathHandling",
			)
			if diff := cmp.Diff(r, tc.expected, opts); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillServiceDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err := FillEntityDefaults(s, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			opt := []cmp.Option{
				cmpopts.IgnoreFields(Service{}, "Enabled"),
			}
			if diff := cmp.Diff(s, tc.expected, opt...); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillTargetDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			target := tc.target
			fullSchema, err := client.Schemas.Get(defaultCtx, "targets")
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err := FillEntityDefaults(target, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			if diff := cmp.Diff(target, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillUpstreamsDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			u := tc.upstream
			fullSchema, err := client.Schemas.Get(defaultCtx, "upstreams")
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err = FillEntityDefaults(u, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			// Ignore fields to make tests pass despite small differences across releases.
			opts := []cmp.Option{
				cmpopts.IgnoreFields(Healthcheck{}, "Threshold"),
				cmpopts.IgnoreFields(Upstream{}, "UseSrvName"),
			}
			if diff := cmp.Diff(u, tc.expected, opts...); diff != "" {
				t.Errorf(diff)
			}
		})
	}
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
				t.Errorf(diff)
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
				t.Errorf(diff)
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
				t.Errorf(diff)
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	assert.NoError(err)
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
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err := FillEntityDefaults(plugin, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			if diff := cmp.Diff(plugin, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func Test_fillConfigRecord(t *testing.T) {
	tests := []struct {
		name     string
		schema   gjson.Result
		config   Configuration
		expected Configuration
	}{
		{
			name: "fills defaults for all missing fields",
			schema: gjson.Parse(`{
				"fields": {
					"config":
						{
							"type": "record",
							"fields":[
								{
									"enabled":{
										"type":"boolean",
										"default":true,
										"required":true
									}
								},
								{
									"mappings":{
										"required":false,
										"type":"array",
										"elements":{
											"type":"record",
											"fields":[
												{"name":{"type":"string","required":false}},
												{"nationality":{"type":"string","required":false}}
											]
										}
									}
								}
							]
						}
					}
				}`),
			config: Configuration{
				"mappings": []interface{}{
					map[string]interface{}{
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			configSchema, err := getConfigSchema(tc.schema)
			require.NoError(t, err)
			config := fillConfigRecord(configSchema, tc.config)
			require.NotNil(t, config)
			if diff := cmp.Diff(config, tc.expected); diff != "" {
				t.Errorf(diff)
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
					"metrics": defaultMetrics,
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

			assert.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin,
				"Protocols", "Enabled",
			)
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf(diff)
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
			assert.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf(diff)
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
					"bootstrap_servers": []any{
						map[string]any{
							"host": "127.0.0.1",
							"port": float64(42),
						},
					},
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

			assert.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin,
				"Protocols", "Enabled",
			)
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf(diff)
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
						"kong": map[string]any{},
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
			assert.NoError(t, FillPluginsDefaults(plugin, fullSchema))
			opts := cmpopts.IgnoreFields(*plugin, "Enabled", "Protocols")
			if diff := cmp.Diff(plugin, tc.expected, opts); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
