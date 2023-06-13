# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: coprocess_return_overrides.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import coprocess_common_pb2 as coprocess__common__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='coprocess_return_overrides.proto',
  package='coprocess',
  syntax='proto3',
  serialized_options=b'Z\n/coprocess',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n coprocess_return_overrides.proto\x12\tcoprocess\x1a\x16\x63oprocess_common.proto\"\x93\x01\n\x0fReturnOverrides\x12\x15\n\rresponse_code\x18\x01 \x01(\x05\x12\x16\n\x0eresponse_error\x18\x02 \x01(\t\x12\"\n\x07headers\x18\x03 \x03(\x0b\x32\x11.coprocess.Header\x12\x16\n\x0eoverride_error\x18\x04 \x01(\x08\x12\x15\n\rresponse_body\x18\x05 \x01(\tB\x0cZ\n/coprocessb\x06proto3'
  ,
  dependencies=[coprocess__common__pb2.DESCRIPTOR,])




_RETURNOVERRIDES = _descriptor.Descriptor(
  name='ReturnOverrides',
  full_name='coprocess.ReturnOverrides',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='response_code', full_name='coprocess.ReturnOverrides.response_code', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='response_error', full_name='coprocess.ReturnOverrides.response_error', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='headers', full_name='coprocess.ReturnOverrides.headers', index=2,
      number=3, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='override_error', full_name='coprocess.ReturnOverrides.override_error', index=3,
      number=4, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='response_body', full_name='coprocess.ReturnOverrides.response_body', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=72,
  serialized_end=219,
)

_RETURNOVERRIDES.fields_by_name['headers'].message_type = coprocess__common__pb2._HEADER
DESCRIPTOR.message_types_by_name['ReturnOverrides'] = _RETURNOVERRIDES
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ReturnOverrides = _reflection.GeneratedProtocolMessageType('ReturnOverrides', (_message.Message,), {
  'DESCRIPTOR' : _RETURNOVERRIDES,
  '__module__' : 'coprocess_return_overrides_pb2'
  # @@protoc_insertion_point(class_scope:coprocess.ReturnOverrides)
  })
_sym_db.RegisterMessage(ReturnOverrides)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
