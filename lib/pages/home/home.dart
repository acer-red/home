import 'package:flutter/material.dart';
import 'package:acer_red/env/ui.dart';
import 'package:acer_red/env/config.dart';
import 'package:acer_red/services/http/http.dart';
import './basic_info.dart';
import './safe.dart';
import './feedback.dart';
import './setting.dart';
import './about.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int idx = 0;
  final List<Widget> _pages = [
    BasicInfo(),
    Safe(),
    FeedBack(),
    About(),
    Setting(),
  ];
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          onPressed: () {
            Navigator.pop(context);
          },
          icon: logo(),
        ),
      ),
      body: Row(
        children: <Widget>[
          // Left side - Feature list with a modern look
          Container(
            width: 220,
            decoration: BoxDecoration(
              color: Colors.white,
              boxShadow: [
                BoxShadow(
                  color: Colors.grey..withValues(alpha: 0.3),
                  spreadRadius: 2,
                  blurRadius: 5,
                  offset: Offset(0, 3), // changes position of shadow
                ),
              ],
            ),
            child: ListView(
              padding: EdgeInsets.symmetric(vertical: 20),
              children: <Widget>[
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                    icon: Icon(Icons.info_outline, color: Colors.blueGrey),
                    label: Text(
                    '基础信息',
                    style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    onPressed: () {
                    if (idx != 0) {
                      setState(() {
                      idx = 0;
                      });
                    }
                    },
                    style: TextButton.styleFrom(
                    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                    icon: Icon(Icons.security, color: Colors.blueGrey),
                    label: Text(
                    '安全设置',
                    style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    onPressed: () {
                    if (idx != 1) {
                      setState(() {
                      idx = 1;
                      });
                    }
                    },
                    style: TextButton.styleFrom(
                    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                    icon: Icon(
                    Icons.feedback_outlined,
                    color: Colors.blueGrey,
                    ),
                    label: Text(
                    '反馈',
                    style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    onPressed: () {
                    if (idx != 2) {
                      setState(() {
                      idx = 2;
                      });
                    }
                    },
                    style: TextButton.styleFrom(
                    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                    icon: Icon(Icons.help_outline, color: Colors.blueGrey),
                    label: Text(
                    '关于',
                    style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    onPressed: () {
                    if (idx != 3) {
                      setState(() {
                      idx = 3;
                      });
                    }
                    },
                    style: TextButton.styleFrom(
                    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                    ),
                  ),
                ),
                Divider(
                  color: Colors.grey[300],
                  indent: 20,
                  endIndent: 20,
                ), // Visual separator
                Padding(
                  padding: const EdgeInsets.symmetric(
                  horizontal: 20,
                  vertical: 10,
                  ),
                  child: Text(
                  '其他',
                  style: TextStyle(fontSize: 12, color: Colors.grey[600]),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                    icon: Icon(Icons.settings, color: Colors.blueGrey),
                    label: Text(
                    '设置',
                    style: TextStyle(fontWeight: FontWeight.w500),
                    ),
                    onPressed: () {
                    if (idx != 4) {
                      setState(() {
                      idx = 4;
                      });
                    }
                    },
                    style: TextButton.styleFrom(
                    padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: TextButton.icon(
                  icon: Icon(Icons.exit_to_app),
                  label: Text(
                    '退出',
                    style: TextStyle(
                      fontWeight: FontWeight.w500,),
                  ),
                  onPressed: () => logout(),
                  style: TextButton.styleFrom(
                    padding:
                      EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    alignment: Alignment.centerLeft,
                  ),
                  ),
                ),
              ],
            ),
          ),
          // Right side - Content area
          Expanded(child: _pages[idx]),
        ],
      ),
    );
  }

  logout() {
    final id = Settings().getUID();
    Http().userLogout(RequestPostUserLogout(id: id)).then((onValue) {
      if (onValue.isOK) {
        Settings().setLogin(false);
        if (mounted) {
          Navigator.pop(context);
        }
      } else {
        showMsg(onValue.msg);
      }
    });
  }
}
