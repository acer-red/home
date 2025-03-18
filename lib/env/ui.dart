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

  required String text,
}) {
  return TextButton(
    style: ButtonStyle(
      backgroundColor: WidgetStateProperty.all(
        Theme.of(context).colorScheme.primary,
      ),
    ),
    onPressed: func,
    child: Text(
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
