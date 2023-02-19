#include "globals.hpp"



String charArrayToString(char arrChar[], int tam)
{
  String s = "";
  for (int i = 0; i < tam; i++)
  {
    s = s + arrChar[i];
  }
  return s;
}

void sleep(int ms)
{
  vTaskDelay(ms / portTICK_PERIOD_MS);
}
