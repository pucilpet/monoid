# generated by datamodel-codegen:
#   filename:  monoid_protocol.json
#   timestamp: 2022-11-04T19:51:10+00:00

from __future__ import annotations

from enum import Enum
from typing import Any, Dict, List, Optional, Union

from pydantic import BaseModel


class MonoidQueryIdentifier(BaseModel):
    schema_name: str
    schema_group: Optional[str] = None
    identifier: str
    identifier_query: Union[str, int]
    json_schema: Dict[str, Any]


class MonoidRecord(BaseModel):
    schema_name: str
    schema_group: Optional[str] = None
    data: Dict[str, Any]


class MonoidSchema(BaseModel):
    name: str
    group: Optional[str] = None
    json_schema: Dict[str, Any]


class MonoidSiloSpec(BaseModel):
    name: Optional[str] = None
    spec: Optional[Dict[str, Any]] = None


class MonoidSchemasMessage(BaseModel):
    schemas: List[MonoidSchema]


class Status(Enum):
    SUCCESS = 'SUCCESS'
    FAILURE = 'FAILURE'


class MonoidValidateMessage(BaseModel):
    status: Status
    message: Optional[str] = None


class Type(Enum):
    SCHEMA = 'SCHEMA'
    RECORD = 'RECORD'
    SPEC = 'SPEC'
    VALIDATE = 'VALIDATE'


class MonoidMessage(BaseModel):
    type: Type
    record: Optional[MonoidRecord] = None
    schema_msg: Optional[MonoidSchemasMessage] = None
    spec: Optional[MonoidSiloSpec] = None
    validate_msg: Optional[MonoidValidateMessage] = None


class MonoidProtocol(BaseModel):
    MonoidMessage: Optional[MonoidMessage] = None


class MonoidQuery(BaseModel):
    identifiers: Optional[List[MonoidQueryIdentifier]] = None
