local sys = import '@lib/sys.libsonnet';

{
  hoge: sys.env('HOGE', 'fuga'),
}
