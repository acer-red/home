import 'package:flutter/material.dart';
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
          FittedBox(
            fit: BoxFit.scaleDown,
            child: SelectableText(
              "深度交流",
              style: TextStyle(fontSize: titlesize),
            ),
          ),
          SizedBox(height: 10.0),
          SelectableText("构思中", style: TextStyle(fontSize: 20.0)),
        ],
      ),
    );
  }
}

