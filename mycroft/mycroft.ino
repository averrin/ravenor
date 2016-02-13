#include <SPI.h>

enum { REG_SELECT = 8 }; // пин, управляющий защёлкой (SS в терминах SPI)
static uint8_t mask = 1;
#define INPUT_SIZE 30

void setup()
{
  /* Инициализируем шину SPI. Если используется программная реализация,
   * то вы должны сами настроить пины, по которым будет работать SPI.
   */
  SPI.begin();

  pinMode(REG_SELECT, OUTPUT);
  digitalWrite(REG_SELECT, LOW); // выбор ведомого - нашего регистра
  SPI.transfer(0); // очищаем содержимое регистра
  /* Завершаем передачу данных. После этого регистр установит
   * на выводах QA-QH уровни, соответствующие записанным битам.
   */
  digitalWrite(REG_SELECT, HIGH);
  digitalWrite(REG_SELECT, LOW);
  SPI.transfer(mask);
  digitalWrite(REG_SELECT, HIGH);
  Serial.begin(115200);
  Serial.println("reset");
}


void loop()
{
  while (Serial.available() > 0) {
//      String cmd = Serial.readString();
      String command = Serial.readStringUntil(':');

      if (command == "T") {
        String args = Serial.readStringUntil('\n');
        Serial.println("Toggle led: " + args);
        toggle(args.toInt());
        continue;
      }
      if (command == "R") {
        Serial.println("Reset leds");
        String args = Serial.readStringUntil('\n');
        mask = 0;
        toggle(0);
        continue;
      }
      if (command == "L") {
        String led = Serial.readStringUntil(':');
        String val = Serial.readStringUntil('\n');
        Serial.println("Set led: " + led + " to " + val);
        setLed(led.toInt(), val.toInt());
        continue;
      }

  }
}

void setLed(int bitToSet, int val) {
    bitWrite(mask,bitToSet, val);

    digitalWrite(REG_SELECT, LOW);
    SPI.transfer(mask);
    digitalWrite(REG_SELECT, HIGH);
}

void toggle(int bitToSet) {
    bitWrite(mask,bitToSet, (bitRead(mask,bitToSet) == 0 ? 1 : 0));

    digitalWrite(REG_SELECT, LOW);
    SPI.transfer(mask);
    digitalWrite(REG_SELECT, HIGH);
}
