import 'package:acer_red/env/ui.dart';
import 'package:flutter/material.dart';

class BasicInfo extends StatefulWidget {
  const BasicInfo({super.key});

  @override
  State<BasicInfo> createState() => _BasicInfo();
}

class _BasicInfo extends State<BasicInfo> {
  bool isEditMode = false;
  TextEditingController nickName = TextEditingController();
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: null,
      body: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('昵称：'),
              isEditMode
                  ? SizedBox(
                    width: 5 * 24,
                    child: TextField(
                      controller: nickName,
                      decoration: InputDecoration(
                        hintText: nickName.text.isEmpty ? "未设置" : '',
                      ),
                    ),
                  )
                  : Text(nickName.text.isEmpty ? "未设置" : ''),
            ],
          ),
          Padding(
            padding: const EdgeInsets.only(top: 20.0, bottom: 20.0),
            child: Divider(
              height: 20, // 分割线高度 (包含上下间距)
              thickness: 1, // 分割线粗细
              indent: 20, // 左侧缩进
              endIndent: 20, // 右侧缩进
              color: Colors.grey[200], // 分割线颜色
            ),
          ),
          SizedBox(
            height: 40,
            width: 180,
            child:
                isEditMode
                    ? blackTextButton(
                      context,
                      () {
                        setState(() {
                          isEditMode = false;
                        });
                      },
                      text: '保存',
                      icon: Icon(Icons.save),
                      iconColor: Colors.white,
                    )
                    : blackTextButton(
                      context,
                      () {
                        setState(() {
                          isEditMode = true;
                        });
                      },
                      text: '编辑',
                      icon: Icon(Icons.edit),
                      iconColor: Colors.white,
                    ),
          ),
        ],
      ),
    );
  }
}
