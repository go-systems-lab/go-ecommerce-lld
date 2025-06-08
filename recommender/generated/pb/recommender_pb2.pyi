from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ProductInteraction(_message.Message):
    __slots__ = ("user_id", "product_id", "interaction_type")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    PRODUCT_ID_FIELD_NUMBER: _ClassVar[int]
    INTERACTION_TYPE_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    product_id: str
    interaction_type: str
    def __init__(self, user_id: _Optional[str] = ..., product_id: _Optional[str] = ..., interaction_type: _Optional[str] = ...) -> None: ...

class RecommendationRequestForUserId(_message.Message):
    __slots__ = ("user_id", "skip", "take")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    SKIP_FIELD_NUMBER: _ClassVar[int]
    TAKE_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    skip: int
    take: int
    def __init__(self, user_id: _Optional[str] = ..., skip: _Optional[int] = ..., take: _Optional[int] = ...) -> None: ...

class RecommendationRequestOnViews(_message.Message):
    __slots__ = ("ids", "skip", "take")
    IDS_FIELD_NUMBER: _ClassVar[int]
    SKIP_FIELD_NUMBER: _ClassVar[int]
    TAKE_FIELD_NUMBER: _ClassVar[int]
    ids: _containers.RepeatedScalarFieldContainer[str]
    skip: int
    take: int
    def __init__(self, ids: _Optional[_Iterable[str]] = ..., skip: _Optional[int] = ..., take: _Optional[int] = ...) -> None: ...

class ProductReplica(_message.Message):
    __slots__ = ("id", "name", "description", "price")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    description: str
    price: float
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ..., description: _Optional[str] = ..., price: _Optional[float] = ...) -> None: ...

class RecommendationResponse(_message.Message):
    __slots__ = ("recommended_products",)
    RECOMMENDED_PRODUCTS_FIELD_NUMBER: _ClassVar[int]
    recommended_products: _containers.RepeatedCompositeFieldContainer[ProductReplica]
    def __init__(self, recommended_products: _Optional[_Iterable[_Union[ProductReplica, _Mapping]]] = ...) -> None: ...
