import 'package:acer_red/pages/login/login.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:acer_red/pages/product/entanglement.dart';
import 'package:acer_red/pages/product/feng_xin.dart';
import 'package:acer_red/pages/product/socra_quest.dart';
import 'package:acer_red/pages/product/whispering_time.dart';
import 'package:acer_red/env/config.dart';
import 'package:acer_red/env/ui.dart';
import 'package:acer_red/pages/home/home.dart';
import 'package:acer_red/services/http/base.dart';
import 'package:acer_red/services/http/http.dart';

class Index extends StatefulWidget {
  const Index({super.key});
  // const Index({Key? key, this.userData}) : super(key: key);
  @override
  State<Index> createState() => _Index();
}

class _Index extends State<Index> {
  final ScrollController _scrollController = ScrollController();
  double height = 700;
  final double titlesize = 56.0;
  bool isLogin = false;
  late User user;
  final List<String> productName = ['枫迹', '深度交流', '组织信息', '关系通讯'];
  int shortWidth = 400;

  @override
  void initState() {
    super.initState();
    autologin();
  }

  @override
  void didUpdateWidget(covariant Index oldWidget) {
    super.didUpdateWidget(oldWidget);
    // 当父 Widget 传递的数据发生变化时重新赋值
    autologin();
  }

  autologin() async {
    Http().userAutoLogin().then((onValue) {
      if (onValue.isNotOK) {
        return;
      }

      setState(() {
        user = User(
          username: onValue.username,
          email: onValue.email,
          crtime: onValue.crtime,
          profile: onValue.profile,
        );
        isLogin = true;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        surfaceTintColor: Colors.white,
        leading: MediaQuery.of(context).size.width < shortWidth ? logo() : null,
        actions:
            MediaQuery.of(context).size.width > shortWidth
                ? null
                : [isLogin ? avatarIcon() : loginIcon(), menuIcon()],
        title:
            MediaQuery.of(context).size.width < shortWidth
                ? null
                : Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    logo(),
                    Row(
                      children:
                          productName.map((section) {
                            int index = productName.indexOf(section);
                            return TextButton(
                              onPressed: () {
                                _scrollToSection(index);
                              },
                              child: Text(
                                section,
                                style: TextStyle(color: Colors.black),
                              ),
                            );
                          }).toList(),
                    ),
                    isLogin ? avatarIcon() : loginIcon(),
                  ],
                ),
        backgroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        controller: _scrollController,
        child: Column(
          children: [
            Face(height: height, titlesize: titlesize),
            WhisperingTime(height: height, titlesize: titlesize),
            SocraQuest(height: height, titlesize: titlesize),
            Entanglement(height: height, titlesize: titlesize),
            FengXin(height: height, titlesize: titlesize),
          ],
        ),
      ),
    );
  }

  Widget avatarIcon() {
    final String url = "${HTTPConfig.imageURL}/${user.profile.avatar}";
    return Padding(
      padding: const EdgeInsets.all(2.0),
      child: MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: () {
        enterHomePage();
        },
        child: CircleAvatar(
        backgroundImage: NetworkImage(url),
        radius: 15, // Adjust the radius to make the avatar smaller
        ),
      ),
      ),
    );
  }

  Widget loginIcon() {
    return IconButton(
      icon: Icon(Icons.person),
      onPressed: () {
        isLogin ? enterHomePage() : dialogLogin();
      },
    );
  }

  enterHomePage() async {
    final b = await Navigator.push<bool>(
      context,
      MaterialPageRoute(builder: (context) => HomePage(user)),
    );
    if (b != null) {
      setState(() {
        isLogin = b;
      });
    }
  }

  dialogLogin() async {
    final g = await showDialog<bool>(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          child: IntrinsicHeight(child: SizedBox(width: 250, child: Login())),
        );
      },
    );
    if (g == null) {
      return;
    }
    setState(() {
      isLogin = g;
    });

    await Http().userInfo().then((onValue) {
      if (onValue.isNotOK) {
        return;
      }
      setState(() {
        user = User(
          username: onValue.username,
          email: onValue.email,
          crtime: onValue.crtime,
          profile: onValue.profile,
        );
        isLogin = true;
      });
    });
  }

  Widget menuIcon() {
    return PopupMenuButton<int>(
      color: Colors.white,
      onSelected: (int index) {
        _scrollToSection(index);
      },
      child: const Padding(
        padding: EdgeInsets.all(8.0),
        child: Icon(Icons.menu),
      ),
      itemBuilder: (BuildContext context) {
        return List<PopupMenuEntry<int>>.generate(productName.length, (index) {
          return PopupMenuItem<int>(
            value: index,
            child: Text(productName[index]),
          );
        });
      },
    );
  }

  void _scrollToSection(int index) {
    _scrollController.animateTo(
      height * (index + 1),
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeInOut,
    );
  }
}

class Face extends StatelessWidget {
  const Face({super.key, required this.height, required this.titlesize});

  final double height;
  final double titlesize;

  Future<void> _launchUrl(String url) async {
    final Uri uri = Uri.parse(url);
    if (!await launchUrl(uri)) {
      throw Exception('Could not launch $url');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: double.infinity,
      color: Color.fromARGB(255, 249, 249, 249),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 20, right: 20),
            child: FittedBox(
              fit: BoxFit.scaleDown,
              child: Text(
                "创意产品的汇聚地",
                style: TextStyle(fontSize: titlesize, letterSpacing: 20.0),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ),
          IconButton(
            icon: Image.asset(
              'assets/images/github-mark.png',
              width: 30,
              height: 30,
              fit: BoxFit.scaleDown,
            ),
            onPressed: () {
              _launchUrl('https://github.com/acer-red/home');
            },
          ),
        ],
      ),
    );
  }
}
