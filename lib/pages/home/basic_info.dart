import 'package:acer_red/env/ui.dart';
import 'package:acer_red/env/config.dart';
import 'package:flutter/material.dart';

class BasicInfo extends StatefulWidget {
  final User user;
  const BasicInfo(this.user, {super.key});

  @override
  State<BasicInfo> createState() => _BasicInfo();
}

class _BasicInfo extends State<BasicInfo> {
  bool isEditMode = false;
  TextEditingController nickName = TextEditingController();

  @override
  void initState() {
    super.initState();
    nickName.text = widget.user.profile.nickname;
  }

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
                      decoration: const InputDecoration(
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.black, width: 0.5),
                        borderRadius: BorderRadius.all(Radius.circular(5)),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.black, width: 1),
                        borderRadius: BorderRadius.all(Radius.circular(5)),
                      ),
                      ),
                    ),
                  )
                  : Text(nickName.text),
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
