from google.protobuf import any_pb2 as _any_pb2
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SinkType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    INSERT: _ClassVar[SinkType]
    UPDATE: _ClassVar[SinkType]
    LOG: _ClassVar[SinkType]
INSERT: SinkType
UPDATE: SinkType
LOG: SinkType

class DoItemRequest(_message.Message):
    __slots__ = ("sinkId", "sinkType", "data")
    SINKID_FIELD_NUMBER: _ClassVar[int]
    SINKTYPE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    sinkId: str
    sinkType: SinkType
    data: _any_pb2.Any
    def __init__(self, sinkId: _Optional[str] = ..., sinkType: _Optional[_Union[SinkType, str]] = ..., data: _Optional[_Union[_any_pb2.Any, _Mapping]] = ...) -> None: ...

class DoSinkRequest(_message.Message):
    __slots__ = ("sinkId", "sinkType", "data")
    SINKID_FIELD_NUMBER: _ClassVar[int]
    SINKTYPE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    sinkId: str
    sinkType: SinkType
    data: bytes
    def __init__(self, sinkId: _Optional[str] = ..., sinkType: _Optional[_Union[SinkType, str]] = ..., data: _Optional[bytes] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("status",)
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: bool
    def __init__(self, status: bool = ...) -> None: ...
