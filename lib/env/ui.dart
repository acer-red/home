import 'package:flutter/material.dart';

final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

Image logo() {
  return Image.asset(
    'assets/images/icon-transparent-256.png',
    width: 30,
    height: 30,
    fit: BoxFit.scaleDown,
  );
}

Widget blackTextButton(
  BuildContext context,
  Function() func, {
  Widget? icon,
  Color ? iconColor,
  required String text,
}) {
  return TextButton.icon(
    style: ButtonStyle(
      backgroundColor: WidgetStateProperty.all(
        Theme.of(context).colorScheme.primary,
      ),
    ),
    onPressed: func,
    icon: IconTheme(
      data: IconThemeData(
      color:iconColor ?? Theme.of(context).colorScheme.secondary,
      ),
      child: icon!,
    ),
    label: Text(
      text,
      style: TextStyle(color: Theme.of(context).colorScheme.secondary),
    ),
  );
}

void showMsg(String msg) {
  ScaffoldMessenger.of(navigatorKey.currentContext!).showSnackBar(
    SnackBar(
      content: Row(
        children: [
          Icon(Icons.error, color: Colors.white),
          SizedBox(width: 8),
          Text(msg),
        ],
      ),
      duration: Duration(seconds: 2),
      backgroundColor: Theme.of(navigatorKey.currentContext!).primaryColor,
    ),
  );
}

Divider divider() {
  return Divider(
    height: 20, // 分割线高度 (包含上下间距)
    thickness: 1, // 分割线粗细
    indent: 20, // 左侧缩进
    endIndent: 20, // 右侧缩进
    color: Colors.grey[200], // 分割线颜色
  );
}
