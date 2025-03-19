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
  final User user;
  const HomePage(this.user, {super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();

  int idx = 0;
  late List<Widget> _pages;

  @override
  void initState() {
    super.initState();

    _pages = [
      BasicInfo(widget.user),
      Safe(),
      SizedBox.shrink(),
      FeedBack(),
      Setting(),
      About(),
    ];

    Http().userInfo().then((onValue) {
      if (onValue.isOK) {
      } else {
        showMsg(onValue.msg);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: _scaffoldKey,
      appBar: AppBar(
        leading: IconButton(
          onPressed: () {
            Navigator.pop(context);
          },
          icon: logo(),
        ),
        actions: [menuButton()],
      ),
      endDrawer: Drawer(
        child: Padding(
          padding: const EdgeInsets.only(
            bottom: 20.0,
            top: 20.0,
            left: 10.0,
            right: 10,
          ),
          child: Column(
            children: [
              Expanded(
                child: ListView.builder(
                  itemCount: _pages.length,
                  itemBuilder: (BuildContext context, int index) {
                    String label = '';
                    IconData iconData = Icons.info_outline;
                    VoidCallback? onPressed;
                    bool isHighlighted = idx == index;

                    switch (index) {
                      case 0:
                        label = '基础信息';
                        iconData = Icons.info_outline;
                        onPressed = () {
                          if (idx != index) {
                            setState(() {
                              idx = index;
                            });
                          }
                        };
                        break;
                      case 1:
                        label = '账户安全';
                        iconData = Icons.security;
                        onPressed = () {
                          if (idx != index) {
                            setState(() {
                              idx = index;
                            });
                          }
                        };
                        break;
                      case 2:
                        return Column(
                          children: [
                            Divider(
                              color: Colors.grey[300],
                              indent: 20,
                              endIndent: 20,
                            ),
                          ],
                        );
                      case 3:
                        label = '反馈';
                        iconData = Icons.feedback_outlined;
                        onPressed = () {
                          if (idx != index) {
                            setState(() {
                              idx = index;
                            });
                          }
                        };

                        break;
                      case 4:
                        label = '设置';
                        iconData = Icons.settings;
                        onPressed = () {
                          if (idx != index) {
                            setState(() {
                              idx = index;
                            });
                          }
                        };
                      case 5:
                        label = '关于';
                        iconData = Icons.help_outline;
                        onPressed = () {
                          if (idx != index) {
                            setState(() {
                              idx = index;
                            });
                          }
                        };

                        break;
                    }
                    if (index == 2) {
                      return const SizedBox.shrink();
                    }
                    return Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: TextButton.icon(
                        icon: Icon(iconData, color: Colors.blueGrey),
                        label: Text(
                          label,
                          style: TextStyle(color: Colors.black),
                        ),
                        onPressed: onPressed,
                        style: TextButton.styleFrom(
                          padding: EdgeInsets.symmetric(
                            horizontal: 20,
                            vertical: 10,
                          ),
                          backgroundColor:
                              isHighlighted
                                  ? Colors.grey.shade300
                                  : Colors.white,
                          alignment: Alignment.centerLeft,
                        ),
                      ),
                    );
                  },
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  TextButton.icon(
                    onPressed: () => logout(),
                    icon: Icon(Icons.exit_to_app),
                    label: Text('退出登录'),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),

      body: Column(children: [Expanded(child: _pages[idx])]),
    );
  }

  Widget menuButton() {
    return IconButton(
      icon: Icon(Icons.menu),
      onPressed: () {
        if (_scaffoldKey.currentState?.isDrawerOpen ?? false) {
          _scaffoldKey.currentState?.closeEndDrawer();
        } else {
          _scaffoldKey.currentState?.openEndDrawer();
        }
      },
    );
  }

  logout() {
    Http().userLogout().then((onValue) {
      if (onValue.isOK) {
        if (mounted) {
          Navigator.pop(context);
          Navigator.of(context).pop(false);
        }
      } else {
        showMsg(onValue.msg);
      }
    });
  }
}
