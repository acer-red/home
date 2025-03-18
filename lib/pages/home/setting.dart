import 'package:flutter/material.dart';

class Setting extends StatefulWidget{
  const Setting({super.key});

  @override
  State<Setting> createState() => _Setting();
}
class _Setting extends State<Setting>{
  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: null,
      body: Center(
        child: Text('设置'),
      ),
    );
  }
}