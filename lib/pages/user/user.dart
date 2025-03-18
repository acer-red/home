import 'package:flutter/material.dart';
import 'package:acer_red/env/ui.dart';
import 'package:acer_red/services/http/http.dart';
import 'package:acer_red/env/config.dart';

class Login extends StatefulWidget {
  const Login({super.key});

  @override
  State<Login> createState() => _Login();
}

class _Login extends State<Login> {
  bool _isLogin = true;
  TextEditingController userController = TextEditingController();
  TextEditingController emailController = TextEditingController();
  TextEditingController accountController = TextEditingController();
  TextEditingController passwordController = TextEditingController();
  TextEditingController registerPasswordController = TextEditingController();
  TextEditingController passwordAgainController = TextEditingController();

  bool isPrompt = false;
  String prompt = "";
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(32.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[
          Text(
            _isLogin ? '登录' : '注册',
            style: const TextStyle(fontSize: 36, fontWeight: FontWeight.bold),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 48),
          _isLogin ? loginForm() : registerForm(),
          if (isPrompt)
            Padding(
              padding: const EdgeInsets.only(top: 8.0),
              child: Center(
                child: Text(
                  prompt,
                  maxLines: 1,
                  style: TextStyle(color: Colors.red),
                ),
              ),
            ),
          const SizedBox(height: 24),
          TextButton(
            onPressed: () {
              setState(() {
                _isLogin = !_isLogin;
              });
            },
            child: Text(_isLogin ? '还没有账户？去注册' : '已有账户？去登录'),
          ),
        ],
      ),
    );
  }

  void openPrompt() {
    setState(() {
      isPrompt = true;
    });
  }

  void closePrompt() {
    setState(() {
      isPrompt = false;
    });
  }

  bool switchPrompt(bool s) {
    setState(() {
      isPrompt = s;
    });
    return s;
  }

  Widget loginForm() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
          child: TextField(
            controller: accountController,
            decoration: InputDecoration(labelText: '邮箱/用户名'),
            keyboardType: TextInputType.emailAddress,
          ),
        ),

        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 40.0),
          child: TextField(
            controller: passwordController,
            decoration: InputDecoration(labelText: '密码'),
            obscureText: true,
          ),
        ),
        blackTextButton(context, () => login(), text: '登录'),
      ],
    );
  }

  bool checkAccount() {
    if (accountController.text.isEmpty) {
      prompt = '用户名不能为空';
      return false;
    }
    return true;
  }

  bool checkPassword() {
    if (passwordController.text.isEmpty) {
      prompt = '密码不能为空';
      return false;
    }
    return true;
  }

  bool loginCheck() {
    if (!checkAccount()) {
      return false;
    }
    if (!checkPassword()) {
      return false;
    }
    return true;
  }

  void login() {
    final ok = loginCheck();
    if (!ok) {
      openPrompt();
      return;
    }
    closePrompt();

    final String account = accountController.text;
    final String password = passwordController.text;

    Http()
        .userLogin(RequestPostUserLogin(account: account, password: password))
        .then((value) {
          if (value.isOK) {
            Settings().setLogin(true);
            Settings().setUID(value.id);
            if (mounted) {
              Navigator.of(context).pop();
            }
          } else {
            prompt = value.msg;
            switchPrompt(true);
          }
        });
  }

  Widget registerForm() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
          child: TextField(
            controller: userController,
            decoration: InputDecoration(labelText: '用户名'),
            textInputAction: TextInputAction.next,
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
          child: TextField(
            controller: emailController,
            decoration: InputDecoration(labelText: '邮箱'),
            keyboardType: TextInputType.emailAddress,
            textInputAction: TextInputAction.next,
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
          child: TextField(
            controller: registerPasswordController,
            decoration: InputDecoration(labelText: '密码'),
            obscureText: true,
            textInputAction: TextInputAction.next,
          ),
        ),

        Padding(
          padding: const EdgeInsets.only(top: 8.0, bottom: 40.0),
          child: TextField(
            controller: passwordAgainController,
            decoration: InputDecoration(labelText: '确认密码'),
            obscureText: true,
            textInputAction: TextInputAction.done,
            onSubmitted: (_) => register(),
          ),
        ),

        blackTextButton(context, () => register(), text: '注册'),
      ],
    );
  }

  bool checkRegisterPasswd() {
    final password = registerPasswordController.text;

    if (password.isEmpty) {
      prompt = '密码不能为空';
      return false;
    }

    if (password.length < 8) {
      prompt = '密码长度不能少于 8 个字符';
      return false;
    }

    // 检查是否包含大写字母
    if (!password.contains(RegExp(r'[A-Z]'))) {
      prompt = '密码需要包含至少一个大写字母';
      return false;
    }

    // 检查是否包含小写字母
    if (!password.contains(RegExp(r'[a-z]'))) {
      prompt = '密码需要包含至少一个小写字母';
      return false;
    }

    // 检查是否包含数字
    if (!password.contains(RegExp(r''))) {
      prompt = '密码需要包含至少一个数字';
      return false;
    }

    // 检查是否包含特殊字符 (可以根据需求调整)
    if (!password.contains(RegExp(r'[!@#$%^&*(),.?":{}|<>]'))) {
      prompt = '密码需要包含至少一个特殊字符';
      return false;
    }

    return true;
  }

  bool checkUser() {
    final username = userController.text;
    if (username.isEmpty) {
      prompt = '用户名不能为空';
      return false;
    }
    if (username.length < 3) {
      prompt = '用户名长度不能少于 3 个字符';
      return false;
    }

    if (username.length > 20) {
      prompt = '用户名长度不能超过 20 个字符';
      return false;
    }

    // 允许字母、数字、下划线和点
    final allowedChars = RegExp(r'^[a-zA-Z0-9_.]+$');
    if (!allowedChars.hasMatch(username)) {
      prompt = '用户名只能包含字母、数字、下划线和点';
      return false;
    }

    // 避免使用过于简单的数字组合
    final onlyNumbers = RegExp(r'^[0-9]+$');
    if (onlyNumbers.hasMatch(username) && username.length < 6) {
      prompt = '不允许过于简单的纯数字用户名';
      return false;
    }

    // 避免常见的敏感词
    final forbiddenWords = [
      'admin',
      'test',
      'guest',
      'root',
      'administrator',
      'administrators',
      "superuser",
    ];
    if (forbiddenWords.contains(username.toLowerCase())) {
      prompt = '该用户名已被禁用';
      return false;
    }
    return true;
  }

  bool checkEmail() {
    if (emailController.text.isEmpty) {
      prompt = '邮箱不能为空';
      return false;
    }
    if (!RegExp(
      r'^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$',
    ).hasMatch(emailController.text)) {
      prompt = '邮箱格式不正确';
      return false;
    }
    return true;
  }

  bool registerCheck() {
    if (!checkUser()) {
      return false;
    }
    if (!checkEmail()) {
      return false;
    }
    if (!checkRegisterPasswd()) {
      return false;
    }
    if (registerPasswordController.text != passwordAgainController.text) {
      prompt = '密码不一致';
      switchPrompt(true);
      return false;
    }
    return true;
  }

  void register() {
    final ok = registerCheck();
    if (!ok) {
      openPrompt();
      return;
    }
    closePrompt();

    final String user = userController.text;
    final String email = emailController.text;
    final String password = registerPasswordController.text;

    Http()
        .userRegister(
          RequestPostUserRegister(
            username: user,
            email: email,
            password: password,
          ),
        )
        .then((value) {
          if (value.isOK) {
            Settings().setLogin(true);
            Settings().setUID(value.id);

            if (mounted) {
              Navigator.of(context).pop();
            }
          } else {
            prompt = value.msg;
            switchPrompt(true);
          }
        });
  }
}
