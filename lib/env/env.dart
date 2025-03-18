import 'package:logger/logger.dart';


final serverAddress = "https://acer.red";

var log = Logger(
  printer: PrettyPrinter(
    methodCount: 1, // 设置调用堆栈层级为1
  ),
);

