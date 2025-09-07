package api

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *Pagination
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "default values when no query parameters",
			args: args{
				r: &http.Request{
					URL: &url.URL{},
				},
			},
			want: &Pagination{
				Limit:  DefaultLimit,
				Offset: DefaultOffset,
			},
			wantErr: assert.NoError,
		},
		{
			name: "valid limit and offset parameters",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=20&offset=10",
					},
				},
			},
			want: &Pagination{
				Limit:  20,
				Offset: 10,
			},
			wantErr: assert.NoError,
		},
		{
			name: "only limit parameter provided",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=50",
					},
				},
			},
			want: &Pagination{
				Limit:  50,
				Offset: DefaultOffset,
			},
			wantErr: assert.NoError,
		},
		{
			name: "only offset parameter provided",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "offset=25",
					},
				},
			},
			want: &Pagination{
				Limit:  DefaultLimit,
				Offset: 25,
			},
			wantErr: assert.NoError,
		},
		{
			name: "minimum limit boundary",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=1",
					},
				},
			},
			want: &Pagination{
				Limit:  1,
				Offset: DefaultOffset,
			},
			wantErr: assert.NoError,
		},
		{
			name: "maximum limit boundary",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=100",
					},
				},
			},
			want: &Pagination{
				Limit:  100,
				Offset: DefaultOffset,
			},
			wantErr: assert.NoError,
		},
		{
			name: "minimum offset boundary",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "offset=0",
					},
				},
			},
			want: &Pagination{
				Limit:  DefaultLimit,
				Offset: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid limit - not a number",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=abc",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "invalid offset - not a number",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "offset=xyz",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "limit below minimum",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=0",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "limit above maximum",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=101",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "offset below minimum",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "offset=-1",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "empty string parameters should use defaults",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=&offset=",
					},
				},
			},
			want: &Pagination{
				Limit:  DefaultLimit,
				Offset: DefaultOffset,
			},
			wantErr: assert.NoError,
		},
		{
			name: "negative limit",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=-5",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "float values should fail",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=10.5",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "both parameters invalid",
			args: args{
				r: &http.Request{
					URL: &url.URL{
						RawQuery: "limit=invalid&offset=also_invalid",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := PaginationRequest(tt.args.r)
				if !tt.wantErr(t, err, fmt.Sprintf("PaginationRequest(%v)", tt.args.r)) {
					return
				}
				assert.Equalf(t, tt.want, got, "PaginationRequest(%v)", tt.args.r)
			},
		)
	}
}
