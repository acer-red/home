import 'package:flutter/material.dart';

class WhisperingTime extends StatelessWidget {
  const WhisperingTime({
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
          FittedBox(
            fit: BoxFit.scaleDown,
            child: SelectableText(
              "枫迹",
              style: TextStyle(
                fontSize: titlesize,
                letterSpacing: 20.0, // Adjust the value as needed
              ),
            ),
          ),
          SizedBox(height: 20.0),
          Padding(
            padding: const EdgeInsets.only(left: 20, right: 20),
            child: FittedBox(
              fit: BoxFit.scaleDown,
              child: SelectableText(
                "生活不只是日复一日，更是步步印迹",
                style: TextStyle(fontSize: 20.0),
                maxLines: 1,
              ),
            ),
          ),
          SizedBox(height: 20.0),
          SelectableText("即将开源", style: TextStyle(fontSize: 20.0)),

          // Add more introductory text or widgets here
        ],
      ),
    );
  }
}
