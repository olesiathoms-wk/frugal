// Autogenerated by Frugal Compiler (1.13.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

library variety.src.f_testing_unions;

import 'dart:typed_data' show Uint8List;
import 'package:thrift/thrift.dart';
import 'package:variety/variety.dart' as t_variety;
import 'package:actual_base/actual_base.dart' as t_actual_base;

class TestingUnions implements TBase {
  static final TStruct _STRUCT_DESC = new TStruct("TestingUnions");
  static final TField _AN_ID_FIELD_DESC = new TField("AnID", TType.I64, 1);
  static final TField _A_STRING_FIELD_DESC = new TField("aString", TType.STRING, 2);
  static final TField _SOMEOTHERTHING_FIELD_DESC = new TField("someotherthing", TType.I32, 3);
  static final TField _AN_INT16_FIELD_DESC = new TField("AnInt16", TType.I16, 4);
  static final TField _REQUESTS_FIELD_DESC = new TField("Requests", TType.MAP, 5);

  int _anID;
  static const int ANID = 1;
  String _aString;
  static const int ASTRING = 2;
  int _someotherthing;
  static const int SOMEOTHERTHING = 3;
  int _anInt16;
  static const int ANINT16 = 4;
  Map<int, String> _requests;
  static const int REQUESTS = 5;

  bool __isset_anID = false;
  bool __isset_someotherthing = false;
  bool __isset_anInt16 = false;

  TestingUnions() {
  }

  int get anID => this._anID;

  set anID(int anID) {
    this._anID = anID;
    this.__isset_anID = true;
  }

  bool isSetAnID() => this.__isset_anID;

  unsetAnID() {
    this.__isset_anID = false;
  }

  String get aString => this._aString;

  set aString(String aString) {
    this._aString = aString;
  }

  bool isSetAString() => this.aString != null;

  unsetAString() {
    this.aString = null;
  }

  int get someotherthing => this._someotherthing;

  set someotherthing(int someotherthing) {
    this._someotherthing = someotherthing;
    this.__isset_someotherthing = true;
  }

  bool isSetSomeotherthing() => this.__isset_someotherthing;

  unsetSomeotherthing() {
    this.__isset_someotherthing = false;
  }

  int get anInt16 => this._anInt16;

  set anInt16(int anInt16) {
    this._anInt16 = anInt16;
    this.__isset_anInt16 = true;
  }

  bool isSetAnInt16() => this.__isset_anInt16;

  unsetAnInt16() {
    this.__isset_anInt16 = false;
  }

  Map<int, String> get requests => this._requests;

  set requests(Map<int, String> requests) {
    this._requests = requests;
  }

  bool isSetRequests() => this.requests != null;

  unsetRequests() {
    this.requests = null;
  }

  getFieldValue(int fieldID) {
    switch (fieldID) {
      case ANID:
        return this.anID;
      case ASTRING:
        return this.aString;
      case SOMEOTHERTHING:
        return this.someotherthing;
      case ANINT16:
        return this.anInt16;
      case REQUESTS:
        return this.requests;
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  setFieldValue(int fieldID, Object value) {
    switch(fieldID) {
      case ANID:
        if(value == null) {
          unsetAnID();
        } else {
          this.anID = value;
        }
        break;

      case ASTRING:
        if(value == null) {
          unsetAString();
        } else {
          this.aString = value;
        }
        break;

      case SOMEOTHERTHING:
        if(value == null) {
          unsetSomeotherthing();
        } else {
          this.someotherthing = value;
        }
        break;

      case ANINT16:
        if(value == null) {
          unsetAnInt16();
        } else {
          this.anInt16 = value;
        }
        break;

      case REQUESTS:
        if(value == null) {
          unsetRequests();
        } else {
          this.requests = value;
        }
        break;

      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  // Returns true if the field corresponding to fieldID is set (has been assigned a value) and false otherwise
  bool isSet(int fieldID) {
    switch(fieldID) {
      case ANID:
        return isSetAnID();
      case ASTRING:
        return isSetAString();
      case SOMEOTHERTHING:
        return isSetSomeotherthing();
      case ANINT16:
        return isSetAnInt16();
      case REQUESTS:
        return isSetRequests();
      default:
        throw new ArgumentError("Field $fieldID doesn't exist!");
    }
  }

  read(TProtocol iprot) {
    TField field;
    iprot.readStructBegin();
    while(true) {
      field = iprot.readFieldBegin();
      if(field.type == TType.STOP) {
        break;
      }
      switch(field.id) {
        case ANID:
          if(field.type == TType.I64) {
            anID = iprot.readI64();
            this.__isset_anID = true;
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        case ASTRING:
          if(field.type == TType.STRING) {
            aString = iprot.readString();
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        case SOMEOTHERTHING:
          if(field.type == TType.I32) {
            someotherthing = iprot.readI32();
            this.__isset_someotherthing = true;
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        case ANINT16:
          if(field.type == TType.I16) {
            anInt16 = iprot.readI16();
            this.__isset_anInt16 = true;
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        case REQUESTS:
          if(field.type == TType.MAP) {
            TMap elem46 = iprot.readMapBegin();
            requests = new Map<int, String>();
            for(int elem48 = 0; elem48 < elem46.length; ++elem48) {
              int elem49 = iprot.readI32();
              String elem47 = iprot.readString();
              requests[elem49] = elem47;
            }
            iprot.readMapEnd();
          } else {
            TProtocolUtil.skip(iprot, field.type);
          }
          break;
        default:
          TProtocolUtil.skip(iprot, field.type);
          break;
      }
      iprot.readFieldEnd();
    }
    iprot.readStructEnd();

    // check for required fields of primitive type, which can't be checked in the validate method
    validate();
  }

  write(TProtocol oprot) {
    validate();

    oprot.writeStructBegin(_STRUCT_DESC);
    if(isSetAnID()) {
      oprot.writeFieldBegin(_AN_ID_FIELD_DESC);
      oprot.writeI64(anID);
      oprot.writeFieldEnd();
    }
    if(isSetAString() && this.aString != null) {
      oprot.writeFieldBegin(_A_STRING_FIELD_DESC);
      oprot.writeString(aString);
      oprot.writeFieldEnd();
    }
    if(isSetSomeotherthing()) {
      oprot.writeFieldBegin(_SOMEOTHERTHING_FIELD_DESC);
      oprot.writeI32(someotherthing);
      oprot.writeFieldEnd();
    }
    if(isSetAnInt16()) {
      oprot.writeFieldBegin(_AN_INT16_FIELD_DESC);
      oprot.writeI16(anInt16);
      oprot.writeFieldEnd();
    }
    if(isSetRequests() && this.requests != null) {
      oprot.writeFieldBegin(_REQUESTS_FIELD_DESC);
      oprot.writeMapBegin(new TMap(TType.I32, TType.STRING, requests.length));
      for(var elem50 in requests.keys) {
        oprot.writeI32(elem50);
        oprot.writeString(requests[elem50]);
      }
      oprot.writeMapEnd();
      oprot.writeFieldEnd();
    }
    oprot.writeFieldStop();
    oprot.writeStructEnd();
  }

  String toString() {
    StringBuffer ret = new StringBuffer("TestingUnions(");

    if(isSetAnID()) {
      ret.write("anID:");
      ret.write(this.anID);
    }

    if(isSetAString()) {
      ret.write(", ");
      ret.write("aString:");
      if(this.aString == null) {
        ret.write("null");
      } else {
        ret.write(this.aString);
      }
    }

    if(isSetSomeotherthing()) {
      ret.write(", ");
      ret.write("someotherthing:");
      ret.write(this.someotherthing);
    }

    if(isSetAnInt16()) {
      ret.write(", ");
      ret.write("anInt16:");
      ret.write(this.anInt16);
    }

    if(isSetRequests()) {
      ret.write(", ");
      ret.write("requests:");
      if(this.requests == null) {
        ret.write("null");
      } else {
        ret.write(this.requests);
      }
    }

    ret.write(")");

    return ret.toString();
  }

  validate() {
    // check for required fields
    // check that fields of type enum have valid values
  }
}
