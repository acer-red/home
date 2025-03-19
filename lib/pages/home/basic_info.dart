import 'package:acer_red/env/ui.dart';
import 'package:acer_red/env/config.dart';
import 'package:acer_red/services/http/http.dart';
import 'package:flutter/material.dart';

class BasicInfo extends StatefulWidget {
  final User user;
  const BasicInfo(this.user, {super.key});

  @override
  State<BasicInfo> createState() => _BasicInfo();
}

class _BasicInfo extends State<BasicInfo> {
  late User user;

  bool isEditMode = false;
  TextEditingController nickName = TextEditingController();

  @override
  void initState() {
    super.initState();
    user = widget.user;
    nickName.text = user.profile.nickname;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: null,
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildInfoRow(
              '昵称',
              isEditMode
                  ? SizedBox(
                    width: 200,
                    child: TextField(
                      controller: nickName,
                      autofocus: true,
                      decoration: InputDecoration(
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(8.0),
                          borderSide: BorderSide.none,
                        ),
                        filled: true,
                        fillColor: Colors.grey[200],
                        contentPadding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 8,
                        ),
                      ),
                    ),
                  )
                  : SelectableText(
                    nickName.text,
                    style: const TextStyle(fontSize: 16),
                  ),
            ),
            const SizedBox(height: 16),
            _buildInfoRow(
              '用户名',
              SelectableText(
                widget.user.username,
                style: const TextStyle(fontSize: 16),
              ),
            ),
            const SizedBox(height: 16),
            _buildInfoRow(
              '邮箱',
              SelectableText(
                widget.user.email,
                style: const TextStyle(fontSize: 16),
              ),
            ),
            const SizedBox(height: 24),
            Divider(thickness: 1, color: Colors.grey[300]),
            const SizedBox(height: 24),
            Center(
              child: SizedBox(
                height: 42,
                width: 160,
                child: blackTextButton(
                  context,
                  () {
                    if (isEditMode) {
                      save();
                    } else {
                      setState(() {
                        isEditMode = true;
                      });
                    }
                  },
                  icon: Icon(isEditMode ? Icons.save : Icons.edit),
                  text: isEditMode ? '保存' : '编辑',
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, Widget content) {
    return Row(
      children: [
        SizedBox(
          width: 70,
          child: Text(
            '$label:',
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
        ),
        content,
      ],
    );
  }

  save() async {
    final name = nickName.text;
    if (name.isEmpty) {
      showMsg('昵称不能为空');
      return;
    }
    if (name == user.profile.nickname) {
      setState(() {
        isEditMode = false;
      });
      return;
    }

    RequestPutUserInfo req = RequestPutUserInfo(nickname: nickName.text);
    final g = await Http().userUpadte(req);
    if (g.isNotOK) {
      showMsg(g.msg);
      return;
    }
    setState(() {
      isEditMode = false;
      user.profile.nickname = name;
    });
  }
}
