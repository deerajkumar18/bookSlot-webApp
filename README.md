# bookSlot-webApp
A simple golang web application that allows multiple users to book slot which is achieved with the help of golang's - go routine concurrency concept

HTTP framework :

  Gorilla mux router framework is used for creating a service that hosts locally . 

Design pattern :
  
  Strategy pattern :
    Allows the program to decide dynamically the behaviour of the object . This is achieved by using interfaces with abstract methods . These abstract methods are implemented by the type of the user requests . There are two types of user requests , that are booking and cancellation . We decide the algorithm based on the user input . The interchangeability is achieved by defining seperate types and characters that implements the interface . It eases the thread usage where a central decision logic can be managed . 

Contexts :

  Golang's context along with timeout is used to identify child contexts consuming much time causing deadline errors . 

Tests and Mocks :

  Golang's testing library is used for testing . Currently we have used TBD approach to run tests sequentially . Also mock types are defined and implements the DAO interfaces . Mock allows running unit tests that involves external usage of services (for eg databases) that can be bypassed . 
