# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: mediawatch/enrich/v2/enrich.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database

# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(
    b'\n!mediawatch/enrich/v2/enrich.proto\x12\x14mediawatch.enrich.v2"7\n\rEnrichRequest\x12\x12\n\x04\x62ody\x18\x01 \x01(\tR\x04\x62ody\x12\x12\n\x04lang\x18\x02 \x01(\tR\x04lang"\\\n\x06\x45ntity\x12\x12\n\x04text\x18\x01 \x01(\tR\x04text\x12\x12\n\x04type\x18\x02 \x01(\tR\x04type\x12\x14\n\x05score\x18\x03 \x01(\x01R\x05score\x12\x14\n\x05index\x18\x04 \x03(\x05R\x05index"\xb5\x02\n\x03NLP\x12\x1c\n\tstopwords\x18\x01 \x03(\tR\tstopwords\x12\x1a\n\x08keywords\x18\x02 \x03(\tR\x08keywords\x12\x38\n\x08\x65ntities\x18\x03 \x03(\x0b\x32\x1c.mediawatch.enrich.v2.EntityR\x08\x65ntities\x12\x18\n\x07summary\x18\x04 \x01(\tR\x07summary\x12\x34\n\x06topics\x18\x05 \x03(\x0b\x32\x1c.mediawatch.enrich.v2.EntityR\x06topics\x12\x34\n\x06\x63laims\x18\x06 \x03(\x0b\x32\x1c.mediawatch.enrich.v2.EntityR\x06\x63laims\x12\x34\n\x06quotes\x18\x07 \x03(\x0b\x32\x1c.mediawatch.enrich.v2.EntityR\x06quotes"3\n\x04\x44\x61ta\x12+\n\x03nlp\x18\x01 \x01(\x0b\x32\x19.mediawatch.enrich.v2.NLPR\x03nlp"\x86\x01\n\x0e\x45nrichResponse\x12\x16\n\x06status\x18\x01 \x01(\tR\x06status\x12\x12\n\x04\x63ode\x18\x02 \x01(\x05R\x04\x63ode\x12\x18\n\x07message\x18\x03 \x01(\tR\x07message\x12.\n\x04\x64\x61ta\x18\x04 \x01(\x0b\x32\x1a.mediawatch.enrich.v2.DataR\x04\x64\x61ta2\xcc\x05\n\rEnrichService\x12R\n\x03NLP\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12X\n\tStopWords\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12W\n\x08Keywords\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12W\n\x08\x45ntities\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12V\n\x07Summary\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12U\n\x06Topics\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12U\n\x06Quotes\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x12U\n\x06\x43laims\x12#.mediawatch.enrich.v2.EnrichRequest\x1a$.mediawatch.enrich.v2.EnrichResponse"\x00\x42\x99\x01\n\x18\x63om.mediawatch.enrich.v2B\x0b\x45nrichProtoP\x01\xa2\x02\x03MEX\xaa\x02\x14Mediawatch.Enrich.V2\xca\x02\x14Mediawatch\\Enrich\\V2\xe2\x02 Mediawatch\\Enrich\\V2\\GPBMetadata\xea\x02\x16Mediawatch::Enrich::V2b\x06proto3'
)


_ENRICHREQUEST = DESCRIPTOR.message_types_by_name["EnrichRequest"]
_ENTITY = DESCRIPTOR.message_types_by_name["Entity"]
_NLP = DESCRIPTOR.message_types_by_name["NLP"]
_DATA = DESCRIPTOR.message_types_by_name["Data"]
_ENRICHRESPONSE = DESCRIPTOR.message_types_by_name["EnrichResponse"]
EnrichRequest = _reflection.GeneratedProtocolMessageType(
    "EnrichRequest",
    (_message.Message,),
    {
        "DESCRIPTOR": _ENRICHREQUEST,
        "__module__": "mediawatch.enrich.v2.enrich_pb2"
        # @@protoc_insertion_point(class_scope:mediawatch.enrich.v2.EnrichRequest)
    },
)
_sym_db.RegisterMessage(EnrichRequest)

Entity = _reflection.GeneratedProtocolMessageType(
    "Entity",
    (_message.Message,),
    {
        "DESCRIPTOR": _ENTITY,
        "__module__": "mediawatch.enrich.v2.enrich_pb2"
        # @@protoc_insertion_point(class_scope:mediawatch.enrich.v2.Entity)
    },
)
_sym_db.RegisterMessage(Entity)

NLP = _reflection.GeneratedProtocolMessageType(
    "NLP",
    (_message.Message,),
    {
        "DESCRIPTOR": _NLP,
        "__module__": "mediawatch.enrich.v2.enrich_pb2"
        # @@protoc_insertion_point(class_scope:mediawatch.enrich.v2.NLP)
    },
)
_sym_db.RegisterMessage(NLP)

Data = _reflection.GeneratedProtocolMessageType(
    "Data",
    (_message.Message,),
    {
        "DESCRIPTOR": _DATA,
        "__module__": "mediawatch.enrich.v2.enrich_pb2"
        # @@protoc_insertion_point(class_scope:mediawatch.enrich.v2.Data)
    },
)
_sym_db.RegisterMessage(Data)

EnrichResponse = _reflection.GeneratedProtocolMessageType(
    "EnrichResponse",
    (_message.Message,),
    {
        "DESCRIPTOR": _ENRICHRESPONSE,
        "__module__": "mediawatch.enrich.v2.enrich_pb2"
        # @@protoc_insertion_point(class_scope:mediawatch.enrich.v2.EnrichResponse)
    },
)
_sym_db.RegisterMessage(EnrichResponse)

_ENRICHSERVICE = DESCRIPTOR.services_by_name["EnrichService"]
if _descriptor._USE_C_DESCRIPTORS == False:
    DESCRIPTOR._options = None
    DESCRIPTOR._serialized_options = b"\n\030com.mediawatch.enrich.v2B\013EnrichProtoP\001\242\002\003MEX\252\002\024Mediawatch.Enrich.V2\312\002\024Mediawatch\\Enrich\\V2\342\002 Mediawatch\\Enrich\\V2\\GPBMetadata\352\002\026Mediawatch::Enrich::V2"
    _ENRICHREQUEST._serialized_start = 59
    _ENRICHREQUEST._serialized_end = 114
    _ENTITY._serialized_start = 116
    _ENTITY._serialized_end = 208
    _NLP._serialized_start = 211
    _NLP._serialized_end = 520
    _DATA._serialized_start = 522
    _DATA._serialized_end = 573
    _ENRICHRESPONSE._serialized_start = 576
    _ENRICHRESPONSE._serialized_end = 710
    _ENRICHSERVICE._serialized_start = 713
    _ENRICHSERVICE._serialized_end = 1429
# @@protoc_insertion_point(module_scope)
