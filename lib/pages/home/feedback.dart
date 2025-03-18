import 'package:flutter/material.dart';

class FeedBack extends StatefulWidget{
  const FeedBack({super.key});

  @override
  State<FeedBack> createState() => _FeedBack();
}
class _FeedBack extends State<FeedBack>{
  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: null,
      body: Center(
        child: Text('账户安全'),
      ),
    );
  }
}