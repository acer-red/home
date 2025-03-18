import 'package:acer_red/pages/user/user.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:acer_red/pages/product/entanglement.dart';
import 'package:acer_red/pages/product/feng_xin.dart';
import 'package:acer_red/pages/product/socra_quest.dart';
import 'package:acer_red/pages/product/whispering_time.dart';
import 'package:acer_red/env/config.dart';
import 'package:acer_red/env/ui.dart';
import 'package:acer_red/pages/home/home.dart';

class Home extends StatefulWidget {
  const Home({super.key});

  @override
  State<Home> createState() => _Home();
}

class _Home extends State<Home> {
  final ScrollController _scrollController = ScrollController();
  double height = 700;
  final double titlesize = 56.0;
  final List<String> productName = ['枫迹', '深度交流', '组织信息', '关系通讯'];
  int shortWidth = 400;
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
                : [loginIcon(), menuIcon()],
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
                    loginIcon(),
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

  Widget loginIcon() {
    return IconButton(
      icon: Icon(Icons.person),
      onPressed: () {
        Settings().getLogin() ? enterHomePage() : dialogLogin();
      },
    );
  }

  enterHomePage() {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => HomePage()),
    );
  }

  dialogLogin() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          child: IntrinsicHeight(child: SizedBox(width: 250, child: Login())),
        );
      },
    );
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
