enum Method { get, post, put, delete }

class HTTPConfig {
  static const String serverAddress = String.fromEnvironment(
    'SERVER_ADDRESS',
    defaultValue: 'https://acer.red',
  );
  static  final String imageURL = '$serverAddress/image';
}


class Basic {
  int err;
  String msg;
  Basic({required this.err, required this.msg});

  bool get isNotOK => err != 0;

  bool get isOK => err == 0;
  
}

