import 'dart:convert';

import 'package:acer_red/services/http/base.dart';
import 'package:acer_red/env/env.dart';
import 'package:acer_red/env/config.dart';
import 'package:http/browser_client.dart';
import 'package:http/http.dart' as http;

class RequestPostUserLogin {
  final String account;
  final String password;
  RequestPostUserLogin({required this.account, required this.password});
  Map<String, dynamic> toJson() {
    return {'account': account, 'password': password};
  }
}

class ReponsePostUserLogin extends Basic {
  ReponsePostUserLogin({required super.err, required super.msg});

  factory ReponsePostUserLogin.fromJson(Map<String, dynamic> g) {
    return ReponsePostUserLogin(err: g['err'], msg: g['msg']);
  }
}

class RequestPostUserRegister {
  final String username;
  final String email;
  final String password;

  RequestPostUserRegister({
    required this.username,
    required this.email,
    required this.password,
  });
  Map<String, dynamic> toJson() {
    return {'username': username, 'password': password, 'email': email};
  }
}

class ReponsePostUserRegister extends Basic {
  final String id;

  ReponsePostUserRegister({
    required super.err,
    required super.msg,
    required this.id,
  });
  factory ReponsePostUserRegister.fromJson(Map<String, dynamic> g) {
    return ReponsePostUserRegister(
      err: g['err'],
      msg: g['msg'],
      id: g['data'] != null ? g['data']['id'] : '',
    );
  }
}

class ReponsePostUserLogout extends Basic {
  ReponsePostUserLogout({required super.err, required super.msg});
  factory ReponsePostUserLogout.fromJson(Map<String, dynamic> g) {
    return ReponsePostUserLogout(err: g['err'], msg: g['msg']);
  }
}

class ReponseGetUserInfo extends Basic {
  final String username;
  final String email;
  final String crtime;
  final Profile profile;
  ReponseGetUserInfo({
    required super.err,
    required super.msg,
    required this.username,
    required this.email,
    required this.crtime,
    required this.profile,
  });
  factory ReponseGetUserInfo.fromJson(Map<String, dynamic> g) {
    return ReponseGetUserInfo(
      err: g['err'],
      msg: g['msg'],
      username: g['data'] != null ? g['data']['username'] : '',
      email: g['data'] != null ? g['data']['email'] : '',
      crtime: g['data'] != null ? g['data']['crtime'] : '',
      profile:
          g['data'] != null && g['data']['profile'] != null
              ? Profile.fromJson(g['data']['profile'])
              : Profile(nickname: '', avatar: ''),
    );
  }
}

class ReponsePostUserAutoLogin extends Basic {
  final String username;
  final String email;
  final String crtime;
  final Profile profile;
  ReponsePostUserAutoLogin({
    required super.err,
    required super.msg,
    required this.username,
    required this.email,
    required this.crtime,
    required this.profile,
  });
  factory ReponsePostUserAutoLogin.fromJson(Map<String, dynamic> g) {
    return ReponsePostUserAutoLogin(
      err: g['err'] ?? 0,
      msg: g['msg'] ?? '',
      username: g['data'] != null ? g['data']['username'] : '',
      email: g['data'] != null ? g['data']['email'] : '',
      crtime: g['data'] != null ? g['data']['crtime'] : '',
      profile:
          g['data'] != null && g['data']['profile'] != null
              ? Profile.fromJson(g['data']['profile'])
              : Profile(nickname: '', avatar: ''),
    );
  }
}

class Http {
  Future<T> _handleRequest<T>(
    Method method,
    Uri u,
    Function(Map<String, dynamic>) fromJson, {
    Map<String, dynamic>? data,
    Map<String, String>? headers,
  }) async {
    if (data != null) {
      log.d(data);
    }
    final http.Response response;
    final client = http.Client();
    if (client is BrowserClient) {
      client.withCredentials = true;
    }
    headers = {'Cookie': 'login=0195a9de270c7d98b93754f392b59da9'};
    try {
      switch (method) {
        case Method.get:
          response = await client.get(u, headers: headers);
          break;
        case Method.post:
          response = await client.post(
            u,
            body: jsonEncode(data),
            headers: headers,
          );

          break;
        case Method.put:
          response = await client.put(
            u,
            body: jsonEncode(data),
            headers: headers,
          );
          break;
        case Method.delete:
          response = await client.delete(
            u,
            body: jsonEncode(data),
            headers: headers,
          );
          break;
      }
      if (err(response.statusCode)) {
        return fromJson({'err': 1, 'msg': getMsg(response.statusCode)});
      }
    } catch (e) {
      log.e("请求失败\n${e.toString()}");
      return fromJson({'err': 1, 'msg': '登陆失败，请稍后尝试'});
    }
    try {
      return fromJson(jsonDecode(response.body));
    } catch (e) {
      log.e("解析数据失败 ${e.toString()}\n${response.body}");
      return fromJson({'err': 1, 'msg': '未知错误'});
    }
  }

  bool err(int statusCode) {
    return statusCode >= 400;
  }

  Future<ReponsePostUserLogin> userLogin(RequestPostUserLogin req) async {
    final path = "/api/v1/user/login";
    final uri = Uri.parse(HTTPConfig.serverAddress + path);
    return _handleRequest(
      Method.post,
      uri,
      (g) => ReponsePostUserLogin.fromJson(g),
      data: req.toJson(),
    );
  }

  Future<ReponsePostUserAutoLogin> userAutoLogin() async {
    final path = "/api/v1/user/autologin";
    final uri = Uri.parse(HTTPConfig.serverAddress + path);
    return _handleRequest(
      Method.post,
      uri,
      (g) => ReponsePostUserAutoLogin.fromJson(g),
    );
  }

  Future<ReponsePostUserRegister> userRegister(RequestPostUserRegister req) {
    final path = "/api/v1/user/register";
    final uri = Uri.parse(HTTPConfig.serverAddress + path);

    return _handleRequest(
      Method.post,
      uri,
      (g) => ReponsePostUserRegister.fromJson(g),
      data: req.toJson(),
    );
  }

  Future<ReponsePostUserLogout> userLogout() {
    final path = "/api/v1/user/logout";
    final uri = Uri.parse(HTTPConfig.serverAddress + path);

    return _handleRequest(
      Method.post,
      uri,
      (g) => ReponsePostUserLogout.fromJson(g),
    );
  }

  Future<ReponseGetUserInfo> userInfo() async {
    final path = "/api/v1/user/info";
    final uri = Uri.parse(HTTPConfig.serverAddress + path);

    return _handleRequest(
      Method.get,
      uri,
      (g) => ReponseGetUserInfo.fromJson(g),
    );
  }

  getMsg(int statusCode) {
    final String msg;
    switch (statusCode) {
      case 400:
        msg = "用户名或密码错误";
      case 500:
        msg = "服务器错误";
      default:
        msg = "未知错误，稍后重试";
    }
    return msg;
  }
}
