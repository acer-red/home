import 'package:shared_preferences/shared_preferences.dart';

class Settings {
  static final Settings _instance = Settings._internal();
  SharedPreferences? _prefs;

  factory Settings() {
    return _instance;
  }

  Settings._internal();

  Future<void> init() async {
    _prefs ??= await SharedPreferences.getInstance();
  }

  bool getLogin() {
    bool? b = _prefs?.getBool("login");
    if (b == null) {
      b = false;
      setbool("login", b);
      return b;
    }
    return b;
  }

  String getUID() {
    String? str = _prefs?.getString("uid");
    if (str == null) {
      str = '';
      setString("uid", str);
      return str;
    }
    return str;
  }

  Future<bool> setUID(String str) {
    return setString("uid", str);
  }

  Future<bool> setLogin(bool b) {
    return setbool("login", b);
  }

  Future<bool> setbool(String s, bool b) {
    return _prefs!.setBool(s, b);
  }

  Future<bool> setString(String s, String str) {
    return _prefs!.setString(s, str);
  }
}
