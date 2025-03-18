import 'package:flutter/material.dart';

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
          FittedBox(
            fit: BoxFit.scaleDown,
            child: SelectableText(
              "关系通讯",
              style: TextStyle(fontSize: titlesize, color: Colors.white),
            ),
          ),
          SizedBox(height: 10.0),
          SelectableText(
            "计划中",
            style: TextStyle(fontSize: 20.0, color: Colors.white),
          ),
          // Add more introductory text or widgets here
        ],
      ),
    );
  }
}
