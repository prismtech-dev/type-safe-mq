import base64
from dataclasses import dataclass
import socket
import time
from google.protobuf.message import Message
from typing import Type, TypeVar, Generic, Any

T = TypeVar("T", bound=Message)


@dataclass
class Envelope(Generic[T]):
    payload: T
    origin: str
    timestamp: int

    def __post_init__(self):
        if not isinstance(self.origin, str):
            raise TypeError(f"origin must be str, got {type(self.origin)}")
        if not isinstance(self.timestamp, int):
            raise TypeError(f"timestamp must be int, got {type(self.timestamp)}")

    @classmethod
    def from_json(cls, data: dict[str, Any], proto_cls: Type[T]) -> "Envelope[T]":
        """Parse message received from message queue"""
        try:
            proto = proto_cls()
            proto.ParseFromString(data["payload"])
        except Exception as e:
            raise ValueError("Failed to parse payload as protobuf") from e

        return cls(
            payload=proto,
            origin=data["origin"],
            timestamp=data["timestamp"],
        )

    @classmethod
    def from_bytes(cls, data: dict[bytes, bytes], proto_cls: Type[T]) -> "Envelope[T]":
        """Parse raw byte-keyed message from message queue (e.g., Redis Stream)"""
        try:
            # 1. decode keys to str
            decoded = {k.decode("utf-8"): v for k, v in data.items()}

            # 2. extract fields
            raw_payload = decoded["payload"]
            origin = (
                decoded["origin"].decode("utf-8")
                if isinstance(decoded["origin"], bytes)
                else decoded["origin"]
            )
            timestamp = (
                int(decoded["timestamp"].decode("utf-8"))
                if isinstance(decoded["timestamp"], bytes)
                else int(decoded["timestamp"])
            )

            # 3. parse protobuf
            proto = proto_cls()
            proto.ParseFromString(raw_payload)

            return cls(payload=proto, origin=origin, timestamp=timestamp)

        except KeyError as e:
            raise ValueError(f"Missing field in message: {e}")
        except Exception as e:
            raise ValueError("Failed to parse message from bytes") from e

    @staticmethod
    def pack(payload: T) -> "Envelope[T]":
        if not isinstance(payload, Message):
            raise TypeError(f"Expected protobuf Message, got {type(payload)}")
        return Envelope(
            payload=payload,
            origin=socket.gethostname(),  # current pod name (unique)
            timestamp=int(time.time() * 1000),  # timestamp in millisecond
        )

    def to_dict(self) -> dict[str, Any]:
        """Raw dict with protobuf bytes"""
        return {
            "payload": self.payload.SerializeToString(),
            "origin": self.origin,
            "timestamp": self.timestamp,
        }

    def to_json_safe(self) -> dict[str, Any]:
        """Dict with base64 payload suitable for JSON serialization"""
        return {
            "payload": base64.b64encode(self.payload.SerializeToString()).decode(
                "utf-8"
            ),
            "origin": self.origin,
            "timestamp": self.timestamp,
        }
