

class Profile {
  String nickname;
  String avatar;
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
