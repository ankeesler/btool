FROM alpine
RUN apk add g++ make
RUN wget https://github.com/google/googletest/archive/release-1.8.1.tar.gz -O - | tar xzf -
RUN cd googletest-release-1.8.1 \
    && g++ -v -o /tmp/gtest.o -Igoogletest -Igoogletest/include googletest/src/gtest-all.cc -c \
    && ar rv /usr/lib/libgtest.a /tmp/gtest.o \
    && cp -vr googletest/include/gtest /usr/include \
    && g++ -v -o /tmp/gtest_main.o -Igoogletest -Igoogletest/include googletest/src/gtest_main.cc -c \
    && ar rv /usr/lib/libgtest_main.a /tmp/gtest_main.o \
    && g++ -v -o /tmp/gmock.o -Igoogletest -Igoogletest/include -Igooglemock -Igooglemock/include googlemock/src/gmock-all.cc -c \
    && ar rv /usr/lib/libgmock.a /tmp/gmock.o \
    && g++ -v -o /tmp/gmock_main.o -Igoogletest -Igoogletest/include -Igooglemock -Igooglemock/include googlemock/src/gmock_main.cc -c \
    && ar rv /usr/lib/libgmock_main.a /tmp/gmock_main.o \
    && cp -vr googlemock/include/gmock /usr/include
