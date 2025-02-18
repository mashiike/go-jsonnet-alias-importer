local base = import './base.libsonnet';
local embed = import '@testing/embed.libsonnet';

base + embed{
  var: 1,
}
