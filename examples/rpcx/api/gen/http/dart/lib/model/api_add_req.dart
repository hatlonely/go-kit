part of openapi.api;

class ApiAddReq {
  
  int i1 = null;
  
  int i2 = null;
  ApiAddReq();

  @override
  String toString() {
    return 'ApiAddReq[i1=$i1, i2=$i2, ]';
  }

  ApiAddReq.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    i1 = json['i1'];
    i2 = json['i2'];
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (i1 != null)
      json['i1'] = i1;
    if (i2 != null)
      json['i2'] = i2;
    return json;
  }

  static List<ApiAddReq> listFromJson(List<dynamic> json) {
    return json == null ? List<ApiAddReq>() : json.map((value) => ApiAddReq.fromJson(value)).toList();
  }

  static Map<String, ApiAddReq> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, ApiAddReq>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = ApiAddReq.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of ApiAddReq-objects as value to a dart map
  static Map<String, List<ApiAddReq>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<ApiAddReq>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = ApiAddReq.listFromJson(value);
       });
     }
     return map;
  }
}

