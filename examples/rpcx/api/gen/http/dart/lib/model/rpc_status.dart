part of openapi.api;

class RpcStatus {
  
  int code = null;
  
  String message = null;
  
  List<ProtobufAny> details = [];
  RpcStatus();

  @override
  String toString() {
    return 'RpcStatus[code=$code, message=$message, details=$details, ]';
  }

  RpcStatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    code = json['code'];
    message = json['message'];
    details = (json['details'] == null) ?
      null :
      ProtobufAny.listFromJson(json['details']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (code != null)
      json['code'] = code;
    if (message != null)
      json['message'] = message;
    if (details != null)
      json['details'] = details;
    return json;
  }

  static List<RpcStatus> listFromJson(List<dynamic> json) {
    return json == null ? List<RpcStatus>() : json.map((value) => RpcStatus.fromJson(value)).toList();
  }

  static Map<String, RpcStatus> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, RpcStatus>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = RpcStatus.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of RpcStatus-objects as value to a dart map
  static Map<String, List<RpcStatus>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<RpcStatus>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = RpcStatus.listFromJson(value);
       });
     }
     return map;
  }
}

