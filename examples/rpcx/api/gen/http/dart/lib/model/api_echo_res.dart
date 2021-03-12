part of openapi.api;

class ApiEchoRes {
  
  String message = null;
  ApiEchoRes();

  @override
  String toString() {
    return 'ApiEchoRes[message=$message, ]';
  }

  ApiEchoRes.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    message = json['message'];
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (message != null)
      json['message'] = message;
    return json;
  }

  static List<ApiEchoRes> listFromJson(List<dynamic> json) {
    return json == null ? List<ApiEchoRes>() : json.map((value) => ApiEchoRes.fromJson(value)).toList();
  }

  static Map<String, ApiEchoRes> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, ApiEchoRes>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = ApiEchoRes.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of ApiEchoRes-objects as value to a dart map
  static Map<String, List<ApiEchoRes>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<ApiEchoRes>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = ApiEchoRes.listFromJson(value);
       });
     }
     return map;
  }
}

