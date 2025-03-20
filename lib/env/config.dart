class Avatar{
  String name;
  String url;
  Avatar({required this.name, required this.url});
  factory Avatar.fromJson(Map<String, dynamic> g) {
    return Avatar(name: g['name'], url: g['url']);
  }
}

class Profile {
  String nickname;
  Avatar avatar;
  Profile({required this.nickname, required this.avatar});
  factory Profile.fromJson(Map<String, dynamic> g) {
    return Profile(nickname: g['nickname'], avatar: Avatar.fromJson(g['avatar']));
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
