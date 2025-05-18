import base64
from type_safe_mq import Envelope
from . import mock_pb2


def test_envelope_with_test_payload():
    payload = mock_pb2.MockPayload(
        width=1920,
        height=1080,
        confidence=0.98,
        temperature=36.5,
        is_valid=True,
        label="face",
        image_data=b"\xff\xd8\xff\xe0",
        points=[1, 2, 3],
        tags=[
            mock_pb2.MockPayload.Metadata(key="source", value="camera"),
        ],
        status=mock_pb2.MockPayload.Status.OK,
    )

    # 포장
    env: Envelope[mock_pb2.MockPayload] = Envelope.pack(payload)
    assert env.payload.width == 1920
    assert env.payload.label == "face"

    # dict로 직렬화 후 역직렬화
    as_dict = env.to_dict()
    # print(as_dict)
    parsed = Envelope.from_json(as_dict, mock_pb2.MockPayload)

    assert parsed.payload.width == 1920
    assert parsed.payload.status == mock_pb2.MockPayload.Status.OK
    assert parsed.origin == env.origin

    # base64 버전도 확인
    as_json_safe = env.to_json_safe()
    assert isinstance(as_json_safe["payload"], str)
    # base64 디코딩이 잘 되는지 확인
    raw_bytes = base64.b64decode(as_json_safe["payload"])
    decoded = mock_pb2.MockPayload()
    decoded.ParseFromString(raw_bytes)
    assert decoded.width == 1920
    assert decoded.label == "face"
