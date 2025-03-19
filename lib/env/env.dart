import 'package:logger/logger.dart';

var log = Logger(
  printer: PrettyPrinter(
    methodCount: 1, // 设置调用堆栈层级为1
  ),
);