package attachment

import (
	"testing"
	"tiny-forum/internal/model/request"
)

func TestValidateUploadRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *request.UploadPostFileRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid post_image",
			req:     &request.UploadPostFileRequest{Type: "post_image", PostID: 123},
			wantErr: false,
		},
		{
			name:    "post_image with zero post_id should be allowed for pre-creation upload",
			req:     &request.UploadPostFileRequest{Type: "post_image", PostID: 0},
			wantErr: false,
		},
		{
			name:    "valid comment_file",
			req:     &request.UploadPostFileRequest{Type: "comment_file", ReplyID: 456},
			wantErr: false,
		},
		{
			name:    "comment_file without reply_id",
			req:     &request.UploadPostFileRequest{Type: "comment_file", ReplyID: 0},
			wantErr: true,
		},
		{
			name:    "valid avatar",
			req:     &request.UploadPostFileRequest{Type: "avatar"},
			wantErr: false,
		},
		{
			name:    "valid plugin",
			req:     &request.UploadPostFileRequest{Type: "plugin"},
			wantErr: false,
		},
		{
			name:    "valid post_cover",
			req:     &request.UploadPostFileRequest{Type: "post_cover"},
			wantErr: false,
		},
		{
			name:    "valid topic_cover",
			req:     &request.UploadPostFileRequest{Type: "topic_cover"},
			wantErr: false,
		},
		{
			name:    "valid video",
			req:     &request.UploadPostFileRequest{Type: "video"},
			wantErr: false,
		},
		{
			name:    "invalid type",
			req:     &request.UploadPostFileRequest{Type: "unknown_type"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUploadRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUploadRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
