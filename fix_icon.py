#!/usr/bin/env python3
"""
Fix RT_GROUP_ICON in Wails-built exe.
Wails writes RT_ICON but not RT_GROUP_ICON. This script adds it.
"""

import struct
import os
import shutil
from pefile import PE

def build_group_icon_resource(ico_path):
    """Build RT_GROUP_ICON resource bytes from ICO file."""
    with open(ico_path, 'rb') as f:
        ico = f.read()

    _, _, count = struct.unpack_from('<HHH', ico, 0)

    # Collect image data
    images = []
    for i in range(count):
        off = 6 + i * 16
        w, h = ico[off], ico[off + 1]
        bpp = struct.unpack_from('<H', ico, off + 6)[0]
        sz = struct.unpack_from('<I', ico, off + 8)[0]
        img_off = struct.unpack_from('<I', ico, off + 12)[0]
        images.append({
            'w': w if w < 256 else 0,
            'h': h if h < 256 else 0,
            'bpp': bpp,
            'data': ico[img_off:img_off + sz]
        })

    # Build resource blob: GRPICONDIR + ICONDIRENTRY[] + raw image data
    blob = struct.pack('<HHH', 0, 1, count)

    img_offsets = []
    current = 6 + count * 16
    for img in images:
        blob += struct.pack('<BBBBHHII',
            img['w'], img['h'], 0, 0, 1, img['bpp'],
            len(img['data']), current
        )
        img_offsets.append(current)
        current += len(img['data'])
        current = (current + 7) & ~7  # 8-byte align

    for img in images:
        blob += img['data']
        while len(blob) % 8 != 0:
            blob += b'\x00'

    return blob


def main():
    ico_path = r'D:\project\pc-toolbox\build\windows\icon.ico'
    exe_path = r'D:\project\pc-toolbox\build\bin\pc-toolbox.exe'

    # Backup
    backup = exe_path + '.bak'
    if not os.path.exists(backup):
        shutil.copy2(exe_path, backup)
        print(f"Backed up original to {backup}")

    pe = PE(exe_path)

    # Get all existing resources except RT_ICON (3) and RT_GROUP_ICON (14)
    resources_to_keep = []
    if hasattr(pe, 'DIRECTORY_ENTRY_RESOURCE'):
        for entry in pe.DIRECTORY_ENTRY_RESOURCE.entries:
            type_id = entry.id
            if type_id in (3, 14):
                continue  # Skip icon-related, we'll re-add
            # Copy this resource
            try:
                data_entries = []
                if hasattr(entry, 'directory'):
                    for sub in entry.directory.entries:
                        if hasattr(sub, 'directory'):
                            for de in sub.directory.entries:
                                rva = de.data.struct.OffsetToData
                                size = de.data.struct.Size
                                data = bytes(pe.get_data(rva, size))
                                lang = getattr(de.data, 'lang', 0)
                                data_entries.append((lang, data))
                        else:
                            rva = sub.data.struct.OffsetToData
                            size = sub.data.struct.Size
                            data = bytes(pe.get_data(rva, size))
                            lang = getattr(sub.data, 'lang', 0)
                            data_entries.append((lang, data))
                resources_to_keep.append({
                    'type': type_id,
                    'id': getattr(entry, 'id', 0),
                    'data_entries': data_entries
                })
            except:
                pass

    # Build new RT_GROUP_ICON resource data
    group_icon_blob = build_group_icon_resource(ico_path)
    print(f"RT_GROUP_ICON resource: {len(group_icon_blob)} bytes from {ico_path}")

    # Write patched exe using pe.write()
    # The 'write' method doesn't have a 'resources' parameter, so we need a different approach
    # Let's directly modify the binary data

    # Strategy: build a complete new resource section binary blob
    # and patch it into the PE file

    # For now, use the simplest approach: the go-winres generated syso is

    # Actually, let's use a totally different tactic: write a Go program that
    # reads the ICO and replaces the PE resources
    # But for now, just try the pefile approach

    # Let me just try passing resources to pe.write()
    # Check if it supports this kwarg
    try:
        pe.write(exe_path, resources=[(14, 1, 0, group_icon_blob)])
        print(f"Patched: {exe_path}")
    except TypeError as e:
        print(f"pe.write doesn't support resources kwarg: {e}")
        print("Trying manual binary patch...")
        # Fall back to manual approach
        # Build smallest possible resource section
        pe.close()
        print("ERROR: Cannot patch with this pefile version")
        print("Please try: go build + go-winres approach")


if __name__ == '__main__':
    main()
