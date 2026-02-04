import socket

# 1. SETUP: TARGET MAC (STM32)
# Must match the "byte mac[]" in your Arduino code
target_mac = "DEADBEEFFEED"

# 2. SETUP: YOUR MAC's BRIDGE IP
# From your ifconfig: bridge100 -> inet 192.168.2.1
interface_ip = "192.168.2.1" 

# 3. Create Magic Packet
bytes_mac = bytes.fromhex(target_mac)
magic_packet = b'\xff' * 6 + bytes_mac * 16

# 4. Create Socket
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

# 5. Bind to the Bridge Interface
try:
    sock.bind((interface_ip, 0))
    print(f"Bound to Internet Sharing Bridge: {interface_ip}")
except Exception as e:
    print(f"Error binding to IP: {e}")
    exit()

# 6. Send to Broadcast
print(f"Sending Magic Packet for {target_mac}...")
sock.sendto(magic_packet, ("192.168.2.255", 9)) 
print("Packet Sent.")