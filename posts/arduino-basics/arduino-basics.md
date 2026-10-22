---
title: "A note on the basics of Arduino"
created_at: 2019-05-12
description: "A collection of notes I took for understanding the basic I/O principles of the Arduino microcontrollers."
image: ./resources/arduino.jpg
tags:
    - archive
    - electronics
    - iot
---

These notes focus on the ATmega 328p. References below.

[Arduino Basics 101](https://youtu.be/BtLwoNJ6klE)
[Arduino Basics 102](https://youtu.be/YT3birSKLLU)
[Arduino Basics 103](https://youtu.be/EVm0qVJ56II)

[Programming the Arduino in C](https://balau82.wordpress.com/arduino-in-c/)

[ATMega328p Datasheet](https://www.dropbox.com/s/mx4gheqq9hqgkdk/ATMega328p_datasheet.pdf?dl=0)

[AVR Guide](https://sites.google.com/site/qeewiki/books/avr-guide)

---

I use a lot of the concepts I discuss in my on going buck-boost converter project. You can check it out [here](https://github.com/isaacgr/Buck-BoostConverter).

## Main Functions

These are the primary functions that you will be using when interacting with the Arduino.

- digital output
- digital input
- pwm output (analog output)
- analog input

## Basic Specs

Basic specs of an ATMega 328p.

- 8-bit PWM
- 10-bit ADC
- internal pull-up and pull-down resitors (20k)
- max current draw of 40mA from pins

## Fundamental Commands

- digitalWrite

  - sets a pin HIGH/LOW
- digitalRead

  - reads the digital value of a pin
- analogWrite

  - sets a pin to a digital value of 0-255
- analogRead

  - reads the analog value of a pin
- delay

  - sets a delay in ms

## Port Manipulation

Typically usising the analog read/write or digital read/write funcitons can be quite slow. This is because there is quite a bit of code that needs to run when accessing these functions.

An easy way to work around this is to configure the pin manually. Each port pin consits of 3 register bits.

- `DDxn`
- `PORTxn`
- `PINxn`

The `DDxn` bit in the `DDRx` Register selects the direction of the pin. If `DDxn` is written high, then the pin is configured as an output. If low then it is an input.

The `PORTxn` bit determines the state of the pull-up resistors. If it is written high when the pin is configured as an input, the the pull-up resistor is activated. If it is written low when the pin is an input, the pull-up is turned off.

If `PORTxn` is written high when the pin is configured as an output, then the pin is driven high. If it is writen low when the pin is an output then the pin is driven low.

There are 3 port registers. `PORTB` for digital pin 8-13, `PORTC` which is analog input 0-5 and `PORTD` which is digital pin 0-7.

An example on configuring these registers below

```cpp
DDRD &= ~(1<<DDD2);  // configure PORTD pin 2 as input
PORTD |= (1 << PORTD2);    // turn On the Pull-up
```

First, we set the `DDRD` (Data direction register for port D) and set pin 2 as an input (by setting it low). We enable the pull-up resistor on that pin by setting the `PORTD2` bit to 1.

---

## Timers and Interrupts

[Timers and interrupts](https://www.robotshop.com/community/forum/t/arduino-101-timers-and-interrupts/13072)

### Timers

For more accurate functions, such as PWM, the timer registers can be used.

A timer/counter is a piece of hardware built into the arduino microcontroller. Its like a clock, and can be used to measure time events. The timer can be programmed with some special registers to configure things like the prescaler or the mode.

The Ardunio has 3 timers (`Timer0`, `Timer1`, `Timer2`). Timer0 and Timer2 are 8-bit timers (can count from 0-255), while Timer1 is a 16-bit timer (can count from 0-65535).

#### Timer0

- 8-bit timer
- used for functions such as `delay()`, `millis()` and `micros()`
- if you change the timer0 registers, this may influence the arduino timer functions

#### Timer1

- 16-bit timer
- the `Servo Library` uses timer1 on the arduino uno

#### Timer2

- 8-bit timer
- used by the `tone()` function

#### Timer Register

The timer behaviour can be changed through the timer registers. The most important registers are:

- `TCCRx` the Timer/Counter control register. This is where the prescalar is configured.
- `TCNTx` - Timer/Counter register. The actual timer value is stored here
- `OCRx` - the output compare register
- `ICRx` the input capture register (only available with the 16-bit register)
- `TIMSKx` - the Timer/Counter interrupt mask register. This is to enable or disable the timer interrupts
- `TIFRx` - the Timer/Counter interrupt flag register. This indicates a pending timer interrupt

Different clock sources can be selected for each timer indpendently. To calculate the timer frequency you will need:

1. CPU frequency (16MHz max)
2. max timer counter value (255 or 65535)
3. divide cpu frequency by the prescalar value
4. divide the result through the desired frequency
5. make sure this result is less than the max timer counter value

#### Timer Modes

Timers can be configured in one of two modes

- PWM mode

  - the `OCxy` outputs are used to generate pwm signals
- CTC mode

  - clear timer on compare match
  - when the timer counter reachers the compare match register value, the timer will be cleared and the interrupt flag is set

### Interrupts

An interrupt is an external event that holds the running program and runs a special interrupt service routine (ISR). After the ISR is finished, the running program is continued with the next instruction (machine instruction, not line of code).

Before a pending interrupt will be able to call an ISR the following must be true:

- interrupts must be generally enabled
- the according interrupt mask must be enabled

Interrupts can be enabled/disabled with the `interrupts()` or `nointerrupts()` functions. Interrupts are enabled by default. The masks are enabled or disabled by setting bits in the interrupt mask register (`TIMSKx`).

The arduino functions `attachInterrupt()` and `detachInterrupt()` can only be used for external interrupt pins.

#### Timer Interrupts

A timer can generate different types of interrupts. The register and bit definitions can be found in the datasheet.

#### Timer Overflow

Timer overflow means that the timer has reached its limit value. When a timer overflow interrupt occurs, the timer overflow bit `TOVx` is set in the interrupt flag register `TIFRx`. When the timer overflow interrupt enable bit `TOIEx` is set in the interrupt mask register `TIMSKx`, the timer interrupt service routing will be called `ISR(TIMERx_OVF_vect)`.

#### Output Compare Match

The output compare match will compare the timer value to the register value. When they match, the interrupt occurs. The `OCFxy` flag is set in the `TIFRx` interrupt flag register. When the output compare interrupt enable bit, `OCIExy`, is set in the interrupt mask register `TIMSKx` the compare match interrupt service routine will be called `ISR(TIMER_COMPy_vect)`.

#### Timer Input Capture

The input capture register can capture the Timer/Counter value at a given external event on either the input capture pin `ICPx` or the analog comparator pins. When a timer input capture interrupt occurs, the input capture flag `ICFx` will be set in `TIFRx`. When the interrupt enable bit `ICIExy` is set in `TIMSKx` the service routine will be called `ISR(TIMERx_CAPT_vect)`.

An example of setting up a 30hz timer would be something like this

```cpp
int lcd_timer_init(void)
{
  // 30Hz timer
  TCCR2A |= (1<<WGM21); //CTC mode
  OCR2A = 0XFF; // set the TOP value for the timer (the number it will count up to)
  TIMSK2 |= (1<<OCIE2A); //enable timer a interrupt
  TCCR2B |= (1<<CS22) | (1<<CS21) | (1<<CS20); // 1024 prescaler

  return 0;
}

ISR(TIMER2_COMPA_vect)
{
    // do something in the timer 2 interrupt
}
```

#### PWM and Timer

There is a relationship between the timers and the pwm capable outputs. The datasheet shows that the PWM capable pins have names like `OCRxA, OCRxB, OCRxC` where x is the timer number. The PWM functionality is often shared with other functionality.

The output compare registers are used to generate a PWM or variable frequency otuput on the output compare pins.

---

## ADC

[ADC Basics](https://learn.sparkfun.com/tutorials/analog-to-digital-conversion/what-is-the-adc)
[ADC Registers](http://ee-classes.usc.edu/ee459/library/documents/ADC.pdf)
[Analog Input](https://sites.google.com/site/qeewiki/books/avr-guide/analog-input)

### Analog Input

The ATmega328 has 6 different device pins, but only one ADC. The 6 pins are connected to a multiplexer which allows a single pin at a time to be connected to the ADC.

The ADC is 10-bit, so 1024 bits of resolution. If we are measuring a 5V signal, that gives 5V/1024 = 0.005v/bit accuracy.

### Steps to setup the ADC

#### Setup the reference voltage

The reference voltage source will be the high voltage. This uses `REFS1` and `REFS0` in the `ADMUX` register. This voltage needs to be a minimum of 1.0v and a max of Vcc.

#### Set the ADCMUX source

Because we only have a single ADC but 6 ADC pins to choose from, we need to feed the signal into the ADC using the built in multiplexer. This is done by setting the `MUXn` pins in the `ADMUX` register. If these values while the conversion is running, then the conversion will finish befor ethe change takes.

#### ADC Register control

The ADC is 10bit resolution, however the ATmega is 8bit microcontroller, which means we are always working with 8 bit registers. In order to get 10bit data we nhave to split two registers, `ADCL` and `ADCH`. If you use the full 10bits then `ADCLAR` (which controls how the microcontroller is going to output the data) must be low and you must read the `ADCL` register first, because reading the `ADCH` register causes the ADC to update. If you only want 8 bits, then set `ADCLAR` high and only read from `ADCH`.

#### Conversion Timing

In order to get the best conversion, the ADC clock needs to be 50-200kHz.

ex.

```
Processor Clock Speed / 200k = x  (round up to next division factor)

ex.  Using a 1Mhz Clock
1MHrz / 200k =
1,000,000 / 200,000 = 5
next highest = 8
```

#### Modes of operation

##### Single Conversion Mode

Do one conversion and finish. This is the default setting. This can tie up the processor while doing a conversion, but its simple with multiple inputs because one is triggered at a time.

###### Free Running Mode

When one conversion finsihes, the nex one will start. Set the `ADATE` bit high in the `ADCSRA` register and leave the `ADSTn` bits low in the `ADCSRB` register.
Free running requires extra code since we have to figure out which conversion just finished and switch the MUX at the right time.

##### Interrupts

Whichever mode you choose can be used with or without interrupts. The `ADC_vect` ISR will be triggered if these are enabled.

##### Example

```cpp
// this code continually scans ADC0 for an analog signal, using 16Mhz processor clock

#include <avr/io.h>
#include <stdint.h>       // needed for uint8_t

#include <avr/interrupt.h>

volatile uint8_t ADCvalue;    // Global variable, set to volatile if used with ISR

int main(void)
{

    ADMUX = 0;                // use ADC0
    ADMUX |= (1 << REFS0);    // use AVcc as the reference
    ADMUX |= (1 << ADLAR);    // Right adjust for 8 bit resolution

    ADCSRA |= (1 << ADPS2) | (1 << ADPS1) | (1 << ADPS0); // 128 prescale for 16Mhz
    ADCSRA |= (1 << ADATE);   // Set ADC Auto Trigger Enable

    ADCSRB = 0;               // 0 for free running mode

    ADCSRA |= (1 << ADEN);    // Enable the ADC
    ADCSRA |= (1 << ADIE);    // Enable Interrupts

    ADCSRA |= (1 << ADSC);    // Start the ADC conversion

    sei();

    while (1)
    {
        // main loop
    }
}

ISR(ADC_vect)
{
    ADCvalue = ADCH;          // only need to read the high value for 8 bit
    // REMEMBER: once ADCH is read the ADC will update

    // if you need the value of ADCH in multiple spots, read it into a register
    // and use the register and not the ADCH
}
```

If you need to change the `MUX` in while in free running mode:

```cpp
ISR(ADC_vect)
{
    uint8_t tmp;            // temp register for storage of misc data

    tmp = ADMUX;            // read the value of ADMUX register
    tmp &= 0x0F;            // AND the first 4 bits (value of ADC pin being used)

    ADCvalue = ADCH;        // read the sensor value

    if (tmp == 0)
    {
        // put ADCvalue into whatever register you use for ADC0 sensor
        ADMUX++;            // add 1 to ADMUX to go to the next sensor
    }

    else if (tmp == 1)
    {
        // put ADCvalue into whatever register you use for ADC1 sensor
        ADMUX++;            // add 1 to ADMUX to go to the next sensor
    }
    else if (tmp == 2)
        // put ADCvalue into whatever register you use for ADC2 sensor
        ADMUX &= 0xF8;      // clear the last 4 bits to reset the mux to ADC0
    }
}
```
