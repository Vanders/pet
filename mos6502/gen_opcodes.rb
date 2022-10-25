#!/usr/bin/env ruby

Modes = {
  'IM' => {
    :mode => 'IMMEDIATE',
    :bytes => 1,
    :format => '#$%02x'
  },
  'RE' => {
    :mode => 'RELATIVE',
    :bytes => 1,
    :format => '$%02x'
  },
  'AC' => {
    :mode => 'ACCUMULATOR',
    :bytes => 0,
    :format => ''
  },
  'AB' => {
    :mode => 'ABSOLUTE',
    :bytes => 2,
    :format => '$%04x'
  },
  'ABX' => {
    :mode => 'ABSOLUTE_X',
    :bytes => 2,
    :format => '$%04x,X'
  },
  'ABY' => {
    :mode => 'ABSOLUTE_Y',
    :bytes => 2,
    :format => '$%04x,Y'
  },
  'ZP' => {
    :mode => 'ZERO_PAGE',
    :bytes => 1,
    :format => '$%02x'
  },
  'ZPX' => {
    :mode => 'ZERO_PAGE_X',
    :bytes => 1,
    :format => '$%02x,X'
  },
  'ZPY' => {
    :mode => 'ZERO_PAGE_Y',
    :bytes => 1,
    :format => '$%02x,Y'
  },
  'IN' => {
    :mode => 'INDIRECT',
    :bytes => 2,
    :format => '($%04x)'
  },
  'IX' => {
    :mode => 'INDIRECT_X',
    :bytes => 1,
    :format => '($%02x,X)'
  },
  'IY' => {
    :mode => 'INDIRECT_Y',
    :bytes => 1,
    :format => '($%02x),Y'
  },
  '' => {
    :mode => 'IMPLIED',
    :bytes => 0,
    :format => ''
  }
}

File.readlines('opcodes.go').each do |line|
  if line.include?('INS')
    s = line.split
    ins = s[0]

    s = line.split(/[_ ]/)
    opc = s[1]
    am = s[2]

    m = Modes[am]
    puts "#{ins}: {#{m[:mode]}, #{m[:bytes]}, \"#{opc} #{m[:format]}\", op_#{opc.downcase}},"
  end
end
