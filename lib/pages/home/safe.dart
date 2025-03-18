import 'package:flutter/material.dart';

class Safe extends StatefulWidget{
  const Safe({super.key});

  @override
  State<Safe> createState() => _Safe();
}
class _Safe extends State<Safe>{
  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: null,
      body: Center(
        child: Text('安全设置'),
      ),
    );
  }
}