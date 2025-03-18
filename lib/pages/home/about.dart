import 'package:flutter/material.dart';

class About extends StatefulWidget{
  const About({super.key});

  @override
  State<About> createState() => _About();
}
class _About extends State<About>{
  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: null,
      body: Center(
        child: Text('欢迎回来！'),
      ),
    );
  }
}