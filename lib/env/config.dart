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

class Profile {
  final String nickname;
  final String avatar;
  Profile({required this.nickname, required this.avatar});
  factory Profile.fromJson(Map<String, dynamic> g) {
    return Profile(nickname: g['nickname'], avatar: g['avatar']);
  }
}

class User {
  String username;
  String email;
  String crtime;
  Profile profile;
  User({
    required this.username,
    required this.email,
    required this.crtime,
    required this.profile,
  });
}
