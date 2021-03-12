part of openapi.api;

class ApiAddRes {
  
  int val = null;
  ApiAddRes();

  @override
  String toString() {
    return 'ApiAddRes[val=$val, ]';
  }

  ApiAddRes.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    val = json['val'];
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (val != null)
      json['val'] = val;
    return json;
  }

  static List<ApiAddRes> listFromJson(List<dynamic> json) {
    return json == null ? List<ApiAddRes>() : json.map((value) => ApiAddRes.fromJson(value)).toList();
  }

  static Map<String, ApiAddRes> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, ApiAddRes>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = ApiAddRes.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of ApiAddRes-objects as value to a dart map
  static Map<String, List<ApiAddRes>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<ApiAddRes>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = ApiAddRes.listFromJson(value);
       });
     }
     return map;
  }
}

