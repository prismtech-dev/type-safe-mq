package tests

import (
	"encoding/base64"
	"testing"

	"github.com/prismtech-dev/type-safe-mq/go/envelope"
	testpb "github.com/prismtech-dev/type-safe-mq/go/tests/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestEnvelopeWithMockPayload(t *testing.T) {
	// (1) mock payload 생성
	payload := &testpb.MockPayload{
		Width:       1920,
		Height:      1080,
		Confidence:  0.98,
		Temperature: 36.5,
		IsValid:     true,
		Label:       "face",
		ImageData:   []byte{0xff, 0xd8, 0xff, 0xe0},
		Points:      []int32{1, 2, 3},
		Tags: []*testpb.MockPayload_Metadata{
			{Key: "source", Value: "camera"},
		},
		Status: testpb.MockPayload_OK,
	}

	// (2) Pack
	env := envelope.Pack(payload)
	require.Equal(t, int32(1920), env.Payload.GetWidth())
	require.Equal(t, "face", env.Payload.GetLabel())

	// (3) 직렬화 및 역직렬화
	dict, err := env.ToMap()
	require.NoError(t, err)

	var parsed testpb.MockPayload
	out, err := envelope.FromJSON(dict, &parsed)
	require.NoError(t, err)
	require.Equal(t, payload.GetWidth(), out.Payload.GetWidth())
	require.Equal(t, payload.GetStatus(), out.Payload.GetStatus())
	require.Equal(t, env.Origin, out.Origin)

	// (4) JSON-safe (base64 인코딩) 버전 확인
	jsonSafe, err := env.ToJSONSafe()
	require.NoError(t, err)

	base64Payload, ok := jsonSafe["payload"].(string)
	require.True(t, ok)

	decodedBytes, err := base64.StdEncoding.DecodeString(base64Payload)
	require.NoError(t, err)

	var decoded testpb.MockPayload
	err = proto.Unmarshal(decodedBytes, &decoded)
	require.NoError(t, err)
	require.Equal(t, payload.GetWidth(), decoded.GetWidth())
	require.Equal(t, "face", decoded.GetLabel())
}
