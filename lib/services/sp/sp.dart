import 'package:shared_preferences/shared_preferences.dart';

class SP {
  static final SP _instance = SP._internal();
  SharedPreferences? _prefs;

  factory SP() {
    return _instance;
  }

  SP._internal();

  Future<void> init() async {
    _prefs ??= await SharedPreferences.getInstance();
  }

  // String getUID() {
  //   String? str = _prefs?.getString("uid");
  //   if (str == null) {
  //     str = '';
  //     setString("uid", str);
  //     return str;
  //   }
  //   return str;
  // }

  // Future<bool> setUID(String str) {
  //   return setString("uid", str);
  // }

  Future<bool> setbool(String s, bool b) {
    return _prefs!.setBool(s, b);
  }

  Future<bool> setString(String s, String str) {
    return _prefs!.setString(s, str);
  }
}