import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class Index extends StatefulWidget {
  const Index({super.key});

  @override
  State<Index> createState() => _IndexState();
}

class _IndexState extends State<Index> {
  final ScrollController _scrollController = ScrollController();
  double height = 700;
  final double titlesize = 56.0;
  final List<String> productName = ['枫迹', '深度交流', '组织信息', '关系通讯'];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        surfaceTintColor: Colors.white,

        title: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset(
              'assets/images/icon-transparent-256.png',
              width: 30,
              height: 30,
              fit: BoxFit.scaleDown,
            ),
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
            IconButton(
              icon: Icon(Icons.person),
              onPressed: () {
                // Show a bottom sheet
                showModalBottomSheet(
                  context: context,
                  builder: (BuildContext context) {
                    // Display a message
                    return Container(
                      padding: EdgeInsets.all(20),
                      child: Text('此功能暂未开发，敬请期待！'),
                    );
                  },
                );
                // Close the bottom sheet after 2 seconds
              },
            ),
          ],
        ),
        backgroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        controller: _scrollController,
        child: Column(
          children: [
            Face(height: height, titlesize: titlesize),
            WhiperingTime(height: height, titlesize: titlesize),
            SocraQuest(height: height, titlesize: titlesize),
            Entanglement(height: height, titlesize: titlesize),
            FengXin(height: height, titlesize: titlesize),

            // Second Section: Blocks based on productName, colores, and githubURL
          ],
        ),
      ),
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
          Text(
            "创意产品的汇聚地",
            style: TextStyle(fontSize: titlesize, letterSpacing: 20.0),
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

class WhiperingTime extends StatelessWidget {
  const WhiperingTime({
    super.key,
    required this.height,
    required this.titlesize,
  });

  final double height;
  final double titlesize;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: double.infinity,
      color: Color.fromRGBO(255, 238, 227, 1),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SelectableText(
            "枫迹",
            style: TextStyle(
              fontSize: titlesize,
              letterSpacing: 20.0, // Adjust the value as needed
            ),
          ),
          SizedBox(height: 20.0),
          SelectableText("生活不只是日复一日，更是步步印迹", style: TextStyle(fontSize: 20.0)),
          SizedBox(height: 20.0),
          SelectableText("即将开源", style: TextStyle(fontSize: 20.0)),

          // Add more introductory text or widgets here
        ],
      ),
    );
  }
}

class SocraQuest extends StatelessWidget {
  const SocraQuest({super.key, required this.height, required this.titlesize});

  final double height;
  final double titlesize;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: double.infinity,
      color: Color(0xFFb8b8f3), // Use the second color
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SelectableText("深度交流", style: TextStyle(fontSize: titlesize)),
          SizedBox(height: 10.0),
          SelectableText("构思中", style: TextStyle(fontSize: 20.0)),
        ],
      ),
    );
  }
}

class Entanglement extends StatelessWidget {
  const Entanglement({
    super.key,
    required this.height,
    required this.titlesize,
  });

  final double height;
  final double titlesize;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: double.infinity,
      color: Color(0xFF191923),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SelectableText(
            "组织信息",
            style: TextStyle(fontSize: titlesize, color: Colors.white),
          ),
           SizedBox(height: 10.0),
          SelectableText("计划中", style: TextStyle(fontSize: 20.0,color: Colors.white)),
        ],
      ),
    );
  }
}

class FengXin extends StatelessWidget {
  const FengXin({super.key, required this.height, required this.titlesize});

  final double height;
  final double titlesize;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      width: double.infinity,
      color: Color(0xFFCC5803),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SelectableText(
            "关系通讯",
            style: TextStyle(fontSize: titlesize, color: Colors.white),
          ),
            SizedBox(height: 10.0),
          SelectableText("计划中", style: TextStyle(fontSize: 20.0,color: Colors.white)),
          // Add more introductory text or widgets here
        ],
      ),
    );
  }
}
