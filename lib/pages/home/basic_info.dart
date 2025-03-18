import 'package:flutter/material.dart';

class BasicInfo extends StatefulWidget{
  const BasicInfo({super.key});

  @override
  State<BasicInfo> createState() => _BasicInfo();
}
class _BasicInfo extends State<BasicInfo>{
  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: null,
      body: Center(
        child: Text('基础信息'),
      ),
    );
  }
}