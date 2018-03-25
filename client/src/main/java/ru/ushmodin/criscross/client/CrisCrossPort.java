package ru.ushmodin.criscross.client;

import io.reactivex.Observable;

public interface CrisCrossPort {
    Observable<Void> registration(String username, String email, String password);
    Observable<Void> authorization(String username, String password);
    Observable<String> startGame();
    Observable<?> getGameStatus();
    Observable<?> step(int row, int col);
}
